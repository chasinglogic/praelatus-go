<template>
  <span>
    <router-link v-bind:to="'/projects/' + ticket.project">
      {{ ticket.project }}
    </router-link>
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
    <router-link v-bind:to="'/tickets/' + ticket.key">
      {{ ticket.key }}
    </router-link>
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
