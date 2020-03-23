<template>
  <v-app id="inspire">
    <v-app-bar app fixed clipped-left>
      <v-toolbar-title>{{$t('app.title')}}</v-toolbar-title>
    </v-app-bar>
    <v-content>
      <v-container fluid class="px-0">
        <game-page></game-page>
      </v-container>
    </v-content>
    <v-footer app fixed>
      <span style="margin-left:1em">&copy; Laurent Bernab&eacute; - 2020</span>
    </v-footer>
  </v-app>
</template>

<script>
import GamePage from "./components/GamePage";
import MovesHistory from "./components/MovesHistory";

export default {
  mounted() {
    this.$i18n.locale = navigator.language.substring(0, 2);
  },
  data: () => ({}),
  props: {
    source: String
  },
  components: {
    "game-page": GamePage
  }
};
</script>

<style>
.logo {
  width: 16em;
}
</style>