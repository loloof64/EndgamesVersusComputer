<template>
  <v-container class="white indigo--text" id="root" width="400" height="600">
    <v-row class="mx-2" id="controlsRow">
      <v-col v-for="btn in buttons" :key="btn.tooltipTrKey" cols="3">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn v-on="on" @click="btn.action()" class="blue lighten-3 red--text text-darken-1">
              <v-icon>{{btn.icon}}</v-icon>
            </v-btn>
          </template>
          <span>{{ $t(btn.tooltipTrKey) }}</span>
        </v-tooltip>
      </v-col>
    </v-row>
    <v-row>
      <v-col move cols="4" class="d-flex justify-end">{{$t('history.header.moveNumber')}}</v-col>
      <v-col move cols="4" class="d-flex justify-end">{{$t('history.header.whiteSide')}}</v-col>
      <v-col move cols="4" class="d-flex justify-end">{{$t('history.header.blackSide')}}</v-col>
    </v-row>
    <v-row v-for="(historyLine, index) in history" :key="historyLine.moveNumber">
      <v-col move cols="4" class="d-flex justify-end">{{historyLine.moveNumber}}</v-col>
      <v-col
        move
        cols="4"
        class="d-flex justify-end"
        :class="{ highlight: isSelectedPosition(index, true) }"
        @click="setPosition(index, true)"
      >{{historyLine.white ? historyLine.white.moveFan : ''}}</v-col>
      <v-col
        move
        cols="4"
        class="d-flex justify-end"
        :class="{ highlight: isSelectedPosition(index, false) }"
        @click="setPosition(index, false)"
      >{{historyLine.black? historyLine.black.moveFan : ''}}</v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "MovesHistory",
  props: {
    history: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      buttons: [
        {
          action: this.goFirst,
          tooltipTrKey: "history.buttons.first",
          icon: "mdi-skip-previous"
        },
        {
          action: this.goPrevious,
          tooltipTrKey: "history.buttons.previous",
          icon: "mdi-chevron-left"
        },
        {
          action: this.goNext,
          tooltipTrKey: "history.buttons.next",
          icon: "mdi-chevron-right"
        },
        {
          action: this.goLast,
          tooltipTrKey: "history.buttons.last",
          icon: "mdi-skip-next"
        }
      ],
      selectedPosition: undefined
    };
  },
  methods: {
    setPosition: function(historyIndex, whitePlayer) {
      let valueToEmit;
      const matchingHistoryLine = this.history[historyIndex];
      if (matchingHistoryLine) {
        const isDefined = whitePlayer
          ? matchingHistoryLine["white"]
          : matchingHistoryLine["black"];
        if (isDefined) valueToEmit = { historyIndex, whitePlayer };
      }
      this.$emit("position_requested", valueToEmit);
    },
    goFirst: function() {
      // We don't process it if no history is defined
      if (this.history.length > 0) {
        this.$emit("position_requested", undefined);
      }
    },
    goPrevious: function() {
      if (this.selectedPosition === undefined) return;
      let newSelectedPosition;

      if (
        this.selectedPosition.historyIndex === 0 &&
        this.selectedPosition.whitePlayer === true
      ) {
        newSelectedPosition = undefined;
      } else {
        const weGoIntoPreviousLine =
          this.selectedPosition.whitePlayer === false;
        if (weGoIntoPreviousLine) {
          const noWhiteMove =
            this.history[this.selectedPosition.historyIndex]["white"] ===
            undefined;
          // Did the game started with a black move ?
          if (noWhiteMove) {
            newSelectedPosition = undefined;
          } else {
            newSelectedPosition = {
              ...this.selectedPosition,
              whitePlayer: true
            };
          }
        } else {
          newSelectedPosition = {
            historyIndex: this.selectedPosition.historyIndex - 1,
            whitePlayer: false
          };
        }
      }
      this.$emit("position_requested", newSelectedPosition);
    },
    goNext: function() {
      if (this.history.length === 0) return;

      const firstMoveLineIndex = this.history.findIndex(
        curr => curr.white || curr.black
      );
      const firstMoveLine = this.history[firstMoveLineIndex];
      const whitePlayer = firstMoveLine.white !== undefined;

      let newSelectedPosition;

      if (this.selectedPosition === undefined) {
        newSelectedPosition = { historyIndex: firstMoveLineIndex, whitePlayer };
      } else {
        let alreadyOnLastMove = false;
        const lastMoveLine =
          this.history.length > 0
            ? this.history[this.history.length - 1]
            : undefined;
        const weAreOnLastLine =
          this.selectedPosition.historyIndex === this.history.length - 1;

        if (weAreOnLastLine) {
          if (lastMoveLine.black !== undefined)
            alreadyOnLastMove = this.selectedPosition.whitePlayer === false;
          else alreadyOnLastMove = true;
        }
        if (!alreadyOnLastMove) {
          const weGoIntoNextLine = this.selectedPosition.whitePlayer === false;
          if (weGoIntoNextLine) {
            newSelectedPosition = {
              historyIndex: this.selectedPosition.historyIndex + 1,
              whitePlayer: true
            };
          } else {
            newSelectedPosition = {
              ...this.selectedPosition,
              whitePlayer: false
            };
          }
        } else return;
      }
      this.$emit("position_requested", newSelectedPosition);
    },
    goLast: function() {
      if (this.history.length === 0) return;

      const lastHistoryLineIndex = this.history.length - 1;
      const lastHistoryLine = this.history[lastHistoryLineIndex];
      const newSelectedPosition = {
        historyIndex: lastHistoryLineIndex,
        whitePlayer: lastHistoryLine.black === undefined
      };

      this.$emit("position_requested", newSelectedPosition);
    },
    confirmPositionSet: function(evt) {
      this.selectedPosition = evt;
      this.adjustScroll();
    },
    adjustScroll: function() {
      let lineIndex = this.selectedPosition.historyIndex;
      const rootElement = document.querySelector("#root");
      rootElement.scroll(0, lineIndex * 60);
    },
    isSelectedPosition: function(index, whitePlayer) {
      if (this.selectedPosition === undefined) return false;
      return (
        whitePlayer === this.selectedPosition.whitePlayer &&
        index === this.selectedPosition.historyIndex
      );
    },
    selectLastMove: function() {
      if (this.history.length === 0) return;
      // The last history move may not have been produced yet by the board component.
      // So it is wiser to wait a little.
      setTimeout(() => {
        const historyIndex = this.history.length - 1;
        const whitePlayer =
          this.history[historyIndex]["black"] !== undefined ? false : true;
        this.selectedPosition = { whitePlayer, historyIndex };
      }, 100);
    },
    clearSelection: function() {
      this.selectedPosition = undefined;
    },
    hasData: function() {
      return this.history.length > 0;
    }
  },
  components: {}
};
</script>

<style scoped>
#root {
  color: black;
  display: inline-block;
  height: 600px;
  overflow-y: scroll;
}

.col[move] {
  border-right: 3px solid black;
  border-bottom: 3px solid black;
  font-size: 24px;
}

.highlight {
  background-color: aqua;
  color: brown;
}

#controlsRow {
  position: fixed;
  left: 680px;
  top: 60px;
}
</style>