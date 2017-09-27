<template>
  <span>
    <a v-bind:href="'/#/projects/' + ticket.project">
      {{ ticket.project }}
    </a>
    <span>
      /
    </span>
    <template v-if="ticket.parent" >
      <a v-if="ticket.parent" v-bind:href="'/#/tickets/' + ticket.parent">
        {{ ticket.parent }}
      </a>
      <span>
        /
        </span>
    </template>
    <a v-bind:href="'/#/tickets/' + ticket.key">
      {{ ticket.key }}
    </a>
  </span>
</template>


<script>
 export default {
   name: 'ticket-bread-crumb',
   props: {
     ticket: {}
   }
 }
</script>
