<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div v-bind:style="{ backgroundColor: bg, color: fg, border: brd, padding: pad, fontWeight: 'bold', borderRadius: '5px', width: '100%', textAlign: 'center' }">
    {{ status }}
  </div>
</template>

<script>
 export default {
   name: 'status-pill',

   computed: {
     statusType: function () {
       if (this.status.toLowerCase() === 'backlog' || this.status.toLowerCase() === 'todo') {
         return 'TODO'
       } else if (this.status.toLowerCase() === 'in progress') {
         return 'IN_PROGRESS'
       } else if (this.status.toLowerCase() === 'closed' || this.status.toLowerCase() === 'done') {
         return 'DONE'
       } else {
         return null
       }
     },

     pad: function () {
       return '0.5rem'
     },

     brd: function () {
       if (this.bg === 'white') {
         return 'solid 1px black'
       }
     },

     bg: function () {
       switch (this.statusType) {
         case 'TODO':
           return 'blue'
         case 'IN_PROGRESS':
           return 'yellow'
         case 'DONE':
           return 'green'
         default:
           return 'white'
       }
     },

     fg: function () {
       if (this.bg === 'white') {
         return 'black'
       }

       return 'white'
     }
   },

   props: {
     status: {
       name: 'status',
       default: ''
     }
   }
 }
</script>
