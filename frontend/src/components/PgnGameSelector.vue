<template>
  <simple-modal-dialog
    ref="rootDialog"
    :title="$t('modals.pgnSelection.title')"
    :confirmAction="loadSelectedPgn"
    cancelButton
  >
    <loloof64-chessboard
      ref="previewBoard"
      size="300"
      :white_player_human="false"
      :black_player_human="false"
      :left="10"
      :top="10"
      :reversed="previewBoardReversed"
      :coordinates_visible="true"
    ></loloof64-chessboard>

    <div id="previewPanel">
      <div id="previewControls">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn v-on="on" @click="goPreviousPreview()" class="blue red--text">
              <v-icon>mdi-chevron-left</v-icon>
            </v-btn>
          </template>
          <span>{{ $t('pgn.buttons.previous') }}</span>
        </v-tooltip>

        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn v-on="on" @click="goNextPreview()" class="blue red--text">
              <v-icon>mdi-chevron-right</v-icon>
            </v-btn>
          </template>
          <span>{{ $t('pgn.buttons.next') }}</span>
        </v-tooltip>

        <span>{{selectedPgnIndex + 1}} / {{pgnsArray.length}}</span>
      </div>

      <div id="previewPlayers">{{ previewPgnWhitePlayer }} | {{ previewPgnBlackPlayer }}</div>

      <div id="previewEvent">{{ selectedPgnSite }} {{ previewPgnDate }}</div>
    </div>
  </simple-modal-dialog>
</template>

<script>
import SimpleModalDialog from "./SimpleModalDialog";
import pgnReader from "../libs/pgn";

export default {
  data() {
    return {
      selectedPgnIndex: -1,
      selectedPgnWhite: "",
      selectedPgnBlack: "",
      selectedPgnSite: "",
      selectedPgnDate: "",
      previewBoardReversed: false
    };
  },
  props: {
    pgnsArray: {
      type: Array,
      required: true
    }
  },
  methods: {
    open: function() {
      this.selectedPgnIndex = 0;
      this.$refs["rootDialog"].open();
      setTimeout(() => {
        this.previewPgn();
      }, 50);
    },
    loadSelectedPgn: function() {
      this.$emit("pgnSelected", this.selectedPgnIndex);
    },
    previewPgn: function() {
      if (this.selectedPgnIndex > this.pgnsArray.length) return;

      const selectedPgn = this.pgnsArray[this.selectedPgnIndex];
      const loader = new pgnReader({ pgn: selectedPgn });
      const result = loader.load_pgn();
      const finalPosition = result.finalPosition;

      this.selectedPgnWhite = result.headers.White || "";
      this.selectedPgnBlack = result.headers.Black || "";
      this.selectedPgnSite = result.headers.Site || "";
      this.selectedPgnDate = result.headers.Date || "";

      const previewComponent = this.$refs["previewBoard"];
      previewComponent.newGame(finalPosition);
    },
    goPreviousPreview: function() {
      if (this.selectedPgnIndex > 0) {
        this.selectedPgnIndex -= 1;
        this.previewPgn();
      }
    },
    goNextPreview: function() {
      if (this.selectedPgnIndex < this.pgnsArray.length - 1) {
        this.selectedPgnIndex += 1;
        this.previewPgn();
      }
    }
  },
  computed: {
    previewPgnWhitePlayer: function() {
      return this.selectedPgnWhite.length > 0 ? this.selectedPgnWhite : "???";
    },
    previewPgnBlackPlayer: function() {
      return this.selectedPgnBlack.length > 0 ? this.selectedPgnBlack : "???";
    },
    previewPgnDate: function() {
      return this.selectedPgnDate.length > 0
        ? "(" + this.selectedPgnDate + ")"
        : "";
    }
  },
  components: {
    "simple-modal-dialog": SimpleModalDialog
  }
};
</script>

<style scoped>
#pgnSelectionDialog {
  position: relative;
}

#previewBoard {
  position: absolute;
  left: 10px;
}

#previewPanel {
  position: absolute;
  left: 350px;
  top: 60px;
}

#previewControls > .v-btn {
  margin-left: 10px;
  margin-right: 10px;
}

#previewPlayers {
  margin-left: 20px;
  margin-top: 20px;
}

#previewEvent {
  margin-left: 20px;
  margin-top: 20px;
}
</style>