<template>
  <b-form-fieldset class="mr-auto ml-auto" style="max-width: 70%">
      <b-form-input v-model="query" @keyup="callSearchFunction"
          placeholder="Type to Search" />
  </b-form-fieldset>
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
