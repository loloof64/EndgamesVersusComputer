<template>
  <v-layout align-center class="px-0 mx-6">
      <loloof64-chessboard class="mx-2" size="600"></loloof64-chessboard>
      <move-history :history="history" class="mx-2"></move-history>
  </v-layout>
</template>

<script>
import MovesHistory from './MovesHistory';
export default {
  data() {
    return {
      history: [],
      orderedHistory: []
    };
  },
  components: {
    "move-history": MovesHistory,
  }
};
</script>