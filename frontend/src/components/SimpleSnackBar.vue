<template>
    <v-snackbar ref="snackbar" v-model="snackbarActive" :timeout="1500">{{message}}</v-snackbar>
</template>

<script>
export default {
    name: 'SimpleSnackBar',
    data() {
        return {
            snackbarActive: false,
            message: '',
        };
    },
    methods: {
        open(message) {
            this.message = message;
            this.snackbarActive = true;
        }
    }
}
</script>