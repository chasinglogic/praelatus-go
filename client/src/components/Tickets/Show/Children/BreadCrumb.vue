<!-- Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

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
