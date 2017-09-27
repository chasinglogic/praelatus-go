<template>

</template>

<script>
export default {
  name: 'search-bar',
  data: function () {
    return {
      'query': ''
    }
  },

  methods: {
    callSearchFunction: function () {
      if (this.searchFunction) return this.searchFunction(this.query)
      return null
    }
  },

  props: [
    'searchFunction'
  ]
}
</script>
