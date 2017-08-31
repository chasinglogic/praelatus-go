<template>
  <div>
    <h1>Tickets</h1>
    <search-bar :searchFunction=loadTickets></search-bar>
    <ticket-list showColumnPicker="true"></ticket-list>
  </div>
</template>

<script>
import TicketList from '@/components/Tickets/List'
import SearchBar from '@/components/General/SearchBar'

export default {
  components: {
    SearchBar,
    TicketList
  },

  methods: {
    loadTickets: function (query) {
      console.log('loading...')
      let url = '/api/tickets'

      if (query && query !== '') {
        url += '?q=' + query
      }
      console.log(url, query)

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
