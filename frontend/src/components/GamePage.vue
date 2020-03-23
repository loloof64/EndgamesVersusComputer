<template>
  <v-layout align-center class="px-0 mx-6">
    <loloof64-chessboard
      class="mx-2"
      size="600"
      :white_player_human="whitePlayerHuman"
      :black_player_human="blackPlayerHuman"
    ></loloof64-chessboard>
    <move-history ref="history" :history="history" class="mx-2"></move-history>
  </v-layout>
</template>

<script>
/*
*
    History should be something like (here simplified)
    [
        {
            moveNumber: 1,
            whiteTurn: true,
            moveFan: "e4"
        },
        {
            moveNumber: 1,
            whiteTurn: false,
            moveFan: "Nc6"
        },
        {
            moveNumber: 2,
            whiteTurn: true,
            moveFan: "Nf3"
        },
        {
            moveNumber: 2,
            whiteTurn: false,
            moveFan: "e5"
        },
    ]

    orderedHistory should be something like (here simplified)
    [
        {
          moveNumber: 1,
          white: {
            moveFan: "e4"
          },
          black: {
            moveFan: "Nc6"
          }
        },
        {
          moveNumber: 2,
          white: {
            moveFan: "Nf3"
          },
          black: {
            moveFan: "e5"
          }
        },
      ]
    */

import MovesHistory from "./MovesHistory";
export default {
  data() {
    return {
      history: [],
      orderedHistory: [],
      whitePlayerHuman: true,
      blackPlayerHuman: true
    };
  },
  components: {
    "move-history": MovesHistory
  },
  methods: {
    newGame: function(startPosition, whitePlayerHuman) {
      this.$refs["history"].clearSelection();
      this.history = [];
      this.updateOrderedHistory();

      const boardComponent = document.querySelector("loloof64-chessboard");
      this.whitePlayerHuman = whitePlayerHuman === true;
      this.blackPlayerHuman = whitePlayerHuman === false;
      this.boardReversed = whitePlayerHuman === false;

      setTimeout(() => boardComponent.newGame(startPosition), 100);
    },
    addMoveToHistory: function(event) {
      this.history = [...this.history, event.detail.moveObject];
      this.updateOrderedHistory();

      const boardComponent = document.querySelector("loloof64-chessboard");
      if (!boardComponent.gameIsInProgress()) {
        this.$refs["history"].selectLastMove();
      }
    },
    notifyCheckmate: function(event) {
      const whiteCheckmated = event.detail.whiteTurnBeforeMove;
      const gameEndStatusMsg = this.$i18n.t("game.ended.checkmate", {
        winner: this.$i18n.t(
          whiteCheckmated ? "game.side.white" : "game.side.black"
        )
      });
      this.$emit("snackbar", gameEndStatusMsg);
    },
    notifyStalemate: function() {
      const gameEndStatusMsg = this.$i18n.t("game.ended.stalemate");
      this.$emit("snackbar", gameEndStatusMsg);
    },
    notifyPerpetualDraw: function() {
      const gameEndStatusMsg = this.$i18n.t("game.ended.perpetualDraw");
      this.$emit("snackbar", gameEndStatusMsg);
    },
    notifyMissingMaterialDraw: function() {
      const gameEndStatusMsg = this.$i18n.t("game.ended.missingMaterial");
      this.$emit("snackbar", gameEndStatusMsg);
    },
    notifyFiftyMovesDraw: function() {
      const gameEndStatusMsg = this.$i18n.t("game.ended.fiftyMovesRule");
      this.$emit("snackbar", gameEndStatusMsg);
    },
    makeComputerMove: function() {},
    convertMoveStringToObject: function(moveString) {
      const start = this.convertAlgebraicCellToCoordinates(
        moveString.substring(0, 2)
      );
      const end = this.convertAlgebraicCellToCoordinates(
        moveString.substring(2, 4)
      );

      const board = document.querySelector("loloof64-chessboard");

      if (moveString.length >= 5) {
        const promotion = moveString.substring(4, 5);
        board.playMove({
          startFile: start.file,
          startRank: start.rank,
          endFile: end.file,
          endRank: end.rank,
          promotion
        });
      } else {
        board.playMove({
          startFile: start.file,
          startRank: start.rank,
          endFile: end.file,
          endRank: end.rank
        });
      }
    },
    convertAlgebraicCellToCoordinates: function(cellStr) {
      const file = cellStr.charCodeAt(0) - "a".charCodeAt(0);
      const rank = cellStr.charCodeAt(1) - "1".charCodeAt(0);

      return { file, rank };
    },
    updateOrderedHistory() {
      let currentMoveNumber = undefined;
      let currentHistoryLine = undefined;
      let update = this.history.reduce((acc, curr) => {
        if (currentMoveNumber !== curr.moveNumber) {
          if (currentHistoryLine !== undefined) {
            acc.push(currentHistoryLine);
          }
          currentMoveNumber = curr.moveNumber;
          currentHistoryLine = { moveNumber: currentMoveNumber };
        }
        if (curr.whiteTurn) {
          currentHistoryLine["white"] = curr;
        } else {
          currentHistoryLine["black"] = curr;
        }

        return acc;
      }, []);
      // Don't forget to push the current edited element if any !
      if (currentHistoryLine !== undefined) {
        update.push(currentHistoryLine);
      }

      this.orderedHistory = update;
    },
    setPosition: function(evt) {
      const board = document.querySelector("loloof64-chessboard");
      let success;

      if (evt !== undefined) {
        const whitePlayer = evt.whitePlayer;
        const historyIndex = evt.historyIndex;

        const historyLine = this.orderedHistory[historyIndex];

        let positionObject = whitePlayer
          ? historyLine.white
          : historyLine.black;

        success = board.setPositionAndLastMove({ ...positionObject });
      } else {
        success = board.setPositionAndLastMove();
      }

      if (success) {
        const historyComponent = this.$refs["history"];
        historyComponent.confirmPositionSet(evt);
      }
    },
    toggleSide: function() {
      this.boardReversed = !this.boardReversed;
    },
    playerHasWhite: function() {
      return this.whitePlayerHuman;
    }
  }
};
</script>