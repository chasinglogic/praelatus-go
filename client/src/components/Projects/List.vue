<template>
  <div class="container">
    <h1>Projects</h1>
    <b-form-fieldset class="mr-auto ml-auto">
      <b-form-input v-model="query" @keyup="loadProjects"
        placeholder="Type to Search" />
    </b-form-fieldset>
    <List :items="projects" uid="key"
      prefix="/projects/" :columns="columns"></List>
  </div>
</template>

<script>
 import List from '@/components/General/List'
 import Axios from 'axios'

 export default {
   name: 'projects',
   components: {
     List
   },

   methods: {
     loadProjects: function () {
       let url = '/api/projects'
       let inst = this

       if (this.query && this.query !== '') {
         url += '?q=' + this.query
       }

       Axios.get(url)
            .then((res) => {
              inst.projects = res.data
            })
            .catch((err) => {
              if (err.response.status === 404) {
                inst.projects = []
                return
              }

              // TODO: Visually que the user that there's been an error
              console.log('ERROR', err)
            })
     }
   },

   data: function () {
     return {
       query: '',
       columns: [
         'key',
         'name',
         'lead',
         'createdDate'
       ],
       projects: []
     }
   },

   created: function () {
     this.loadProjects()
   }
 }
</script>
