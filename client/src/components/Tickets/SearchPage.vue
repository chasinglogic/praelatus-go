<template>
  <h1>Searching...</h1>
  <ticket-list tickets={{tickets}}></ticket-list>
</template>

<script>
import TicketList from '@/components/Ticket/List.vue'

export default {
    data: () => {
        return {
            tickets: [],
        }
    }
}
</script>
