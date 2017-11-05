<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div class="status-pill"
    v-bind:style="{ backgroundColor: bg, color: fg }">
    {{ status.name }}
  </div>
</template>

<script>
 export default {
   name: 'status-pill',

   computed: {
     bg: function () {
       switch (this.status.type) {
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


<style>
 .status-pill {
   border-radius: 5px;
   max-width: 600px;
   margin-right: auto;
   margin-left: auto;
   width: 100%;
   text-align: center;
   border: solid 1px black;
   font-weight: bold;
 }
</style>
