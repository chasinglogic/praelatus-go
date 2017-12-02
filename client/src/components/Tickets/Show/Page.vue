<template>
  <ticket :ticketKey="ticketKey" @ticketRetrievalError="handleError"></ticket>
</template>

<script>
 import Show from '@/components/Tickets/Show/Index'

 export default {
   name: 'ticket-show-page',

   components: {
     ticket: Show
   },

   methods: {
     handleError: function (err) {
       console.log(err)
     }
   },

   computed: {
     ticketKey: function () {
       return this.$route.params.key
     }
   }
 }
</script>
