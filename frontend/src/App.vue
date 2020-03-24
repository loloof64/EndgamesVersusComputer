<template>
  <v-app id="inspire">
    <v-app-bar app fixed clipped-left>
      <v-toolbar-title>{{$t('app.title')}}</v-toolbar-title>
      <toolbar-button :text="$t('menu.newGame.tooltip')" :action="newGameRequest">
        <v-icon>mdi-restart</v-icon>
      </toolbar-button>
      <toolbar-button :text="$t('menu.toggleSide.tooltip')" :action="toggleSide">
        <v-icon>mdi-arrow-up-down</v-icon>
      </toolbar-button>
      <toolbar-button :text="$t('menu.stopGame.tooltip')" :action="stopGameRequest">
        <v-icon>mdi-stop-circle</v-icon>
      </toolbar-button>
    </v-app-bar>
    <v-content>
      <v-container fluid class="px-0">
        <game-page ref="gameZone"></game-page>

        <simple-modal-dialog ref="errorDialog" :title="errorDialogTitle">
          <v-card-text>{{errorDialogText}}</v-card-text>
        </simple-modal-dialog>

        <pgn-game-selector
          ref="pgnSelectionDialog"
          id="pgnSelectionDialog"
          :pgnsArray="pgnsArray"
          @pgnSelected="loadSelectedPgn"
        ></pgn-game-selector>

        <simple-modal-dialog
          ref="newGameConfirmation"
          :title="$t('modals.newGame.title')"
          :confirmAction="doStartNewGame"
          cancelButton
        >
          <v-card-text>{{$t('modals.newGame.text')}}</v-card-text>
        </simple-modal-dialog>

        <simple-modal-dialog
          ref="stopGameConfirmation"
          :title="$t('modals.stopGame.title')"
          :confirmAction="doStopGame"
          cancelButton
        >
          <v-card-text>{{$t('modals.stopGame.text')}}</v-card-text>
        </simple-modal-dialog>

        <simple-snack-bar ref="snackbar">{{ snackBarMessage }}</simple-snack-bar>
      </v-container>
    </v-content>
    <v-footer app fixed>
      <span style="margin-left:1em">&copy; Laurent Bernab&eacute; - 2020</span>
    </v-footer>
  </v-app>
</template>

<script>
import GamePage from "./components/GamePage";
import ToolbarButton from "./components/ToolbarButton";
import SimpleModalDialog from "./components/SimpleModalDialog";
import SimpleSnackBar from "./components/SimpleSnackBar";
import PgnGameSelector from "./components/PgnGameSelector";

import pgnReader from "./libs/pgn";

export default {
  mounted() {
    this.$i18n.locale = navigator.language.substring(0, 2);
  },
  data: () => ({
    errorDialogTitle: "",
    errorDialogText: "",

    snackBarMessage: "",
    pgnsArray: []
  }),
  props: {
    source: String
  },
  methods: {
    newGameRequest: function() {
      const chessBoard = document.querySelector("loloof64-chessboard");
      if (chessBoard.gameIsInProgress()) {
        this.$refs["newGameConfirmation"].open();
      } else {
        this.doStartNewGame();
      }
    },
    doStartNewGame: function() {
      // Production mode, use window.backend.TextFileManager.GetTextFileContent()
      window.backend.TextFileManager.GetTextFileContentWithPathProviden(
        "/home/laurent-bernabe/Documents/temp/pgn/lesson01.pgn"
      ).then(content => {
        if (content === "#ErrorReadingFile") {
          this.errorDialogTitle = this.$i18n.t("modals.failedToReadPgn.title");
          this.errorDialogText = this.$i18n.t("modals.failedToReadPgn.text");
          this.$refs["errorDialog"].open();
          return;
        }

        const allPgns = this.splitPgn(content);

        if (allPgns.length === 0) {
          this.errorDialogTitle = this.$i18n.t("modals.noGameInPgn.title");
          this.errorDialogText = this.$i18n.t("modals.noGameInPgn.text");
          this.$refs["errorDialog"].open();
          return;
        }

        this.pgnsArray = allPgns;
        this.letUserSelectPgn();
      });
    },
    stopGameRequest: function() {
      const chessBoard = document.querySelector("loloof64-chessboard");
      if (chessBoard.gameIsInProgress()) {
        this.$refs["stopGameConfirmation"].open();
      }
    },
    doStopGame: function() {
      const chessBoard = document.querySelector("loloof64-chessboard");
      chessBoard.stop();
      this.$refs["snackbar"].open(this.$i18n.t("game.status.stopped"));
      this.$refs["gameZone"].$refs["history"].selectLastMove();
    },
    toggleSide: function() {
      this.$refs["gameZone"].toggleSide();
    },
    splitPgn: function(fileContent) {
      let results = [];

      const isOutOfGame = 0;
      const isInHeader = 1;
      const isInGame = 2;

      let status = isOutOfGame;
      let currentGame = "";

      const isAHeaderLine = l => l.startsWith("[");
      const isBlankLine = l => l.trim().length === 0;

      const lines = fileContent.split(/\r?\n/);
      lines.forEach(l => {
        switch (status) {
          case isOutOfGame:
            if (isAHeaderLine(l)) {
              status = isInHeader;
              currentGame += l + "\n";
            } else if (!isBlankLine(l)) {
              status = isInGame;
              currentGame += l + "\n";
            }

            break;

          case isInHeader:
            if (status === isInGame) {
              results.push(currentGame);
              currentGame = "";
            }
            if (isAHeaderLine(l)) {
              currentGame += l + "\n";
            } else if (isBlankLine(l)) {
              status = isInGame;
              currentGame += "\n";
            } else {
              currentGame += l + "\n";
            }

            break;

          case isInGame:
            if (isAHeaderLine(l)) {
              results.push(currentGame);
              currentGame = "";
              status = isInHeader;
              currentGame += l + "\n";
            } else if (isBlankLine(l)) {
              results.push(currentGame);
              currentGame = "";
              status = isOutOfGame;
            } else {
              currentGame += l + "\n";
            }
        }
      });

      if (currentGame.length > 0) results.push(currentGame);

      return results;
    },
    letUserSelectPgn: function() {
      this.selectedPgnIndex = 0;
      this.$refs["pgnSelectionDialog"].open();
    },
    loadSelectedPgn: function(index) {
      const selectedPgn = this.pgnsArray[index];
      const loader = new pgnReader({ pgn: selectedPgn });
      const result = loader.load_pgn();
      const finalPosition = result.finalPosition;

      try {
        const playerHasWhite = finalPosition.split(" ")[1] === "w";

        this.$refs["gameZone"].newGame(finalPosition, playerHasWhite);
        const chessBoard = document.querySelector("loloof64-chessboard");
        if (chessBoard.gameIsInProgress()) {
          this.$refs["snackbar"].open(this.$i18n.t("game.status.started"));
        } else {
          this.$refs["snackbar"].open(
            this.$i18n.t("game.status.already_finished")
          );
        }
      } catch (err) {
        this.errorDialogTitle = this.$i18n.t("modals.failedToReadPgn.title");
        this.errorDialogText = this.$i18n.t("modals.failedToReadPgn.text");
        this.$refs["errorDialog"].open();
      }
    },
  },

  components: {
    "game-page": GamePage,
    "toolbar-button": ToolbarButton,
    "simple-modal-dialog": SimpleModalDialog,
    "simple-snack-bar": SimpleSnackBar,
    "pgn-game-selector": PgnGameSelector
  }
};
</script>

<style scoped>
</style>