<template>
  <div>
    <h1>Tickets</h1>
    <input type="text" placeholder="Search here..."
      v-model="query" @keyup="loadTickets" />
    <ticket-list></ticket-list>
  </div>
</template>

<script>
import TicketList from '@/components/Tickets/List'

export default {
  components: {
    TicketList
  },

  data: function () {
    return {
      'query': ''
    }
  },

  methods: {
    loadTickets: function () {
      console.log('loading...')
      let url = '/api/tickets'

      if (this.query !== '') {
        url += '?q=' + this.query
      }
      console.log(url, this.query)

      this.$store.dispatch('request', {
        url: url,
        key: 'tickets'
      })
    }
  },

  created: function () {
    this.loadTickets()
  }
}
</script>
