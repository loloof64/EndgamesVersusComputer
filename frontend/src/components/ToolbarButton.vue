<template>
    <v-tooltip bottom>
        <template v-slot:activator="{ on }">
            <v-btn icon dark v-on="on" @click="action()"><slot></slot></v-btn>
        </template>
        <span>{{ text }}</span>
    </v-tooltip>
</template>

<script>
export default {
    props: {
        action: {
            type: Function,
            required: true,
        },
        text: {
            type: String,
            required: true,
        }
    }
}
</script>