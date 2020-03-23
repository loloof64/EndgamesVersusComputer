<template>
  <v-dialog v-model="opened" persistent max-width="650">
    <v-card>
      <v-card-title class="headline">{{title}}</v-card-title>
      <slot></slot>
      <v-card-actions>
        <v-btn v-if="cancelButton"
           color="red darken-1"
          text
          @click="closeWithoutAction()"
        >
        {{$t('modals.global.cancelButton')}}
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
          color="blue darken-1"
          text
          @click="closeWithAction()"
        >{{$t('modals.global.okButton')}}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
    name: 'SimpleDialog',

    props: {
      title: {
        type: String,
        required: true,
      },
      confirmAction: {
        type: Function,
        required: false,
      },
      cancelAction: {
        type: Function,
        required: false,
      },
      cancelButton: {
        type: Boolean,
        default: false,
      }
    },

    data() {
        return {
            opened: false,
        }
    },

    methods: {
        open() {
            this.opened = true;
        },
        closeWithAction() {
          this.opened = false;
          if (this.confirmAction) this.confirmAction();
        },
        closeWithoutAction() {
          this.opened = false;
          if (this.cancelAction) this.cancelAction();
        }
    }
}
</script>