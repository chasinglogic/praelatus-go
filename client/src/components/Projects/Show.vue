<template>
  <div class="container">
    <div v-if="project.key">
      <div class="card">
        <h1 class="card-block">
          <img v-bind:src="project.iconUrl" v-if="project.iconUrl" />
          {{ project.key }} / {{ project.name }}
        </h1>
        <div class="card-block">
          <b-navbar toggleable toggle-breakpoint="md">
            <div class="container">
              <b-nav is-nav-bar>
                <b-nav-item v-bind:to="ticketLink">
                  Tickets
                </b-nav-item>
                <b-nav-item v-bind:to="'/projects/' + project.key + '/settings'">
                  Settings
                </b-nav-item>
              </b-nav>
            </div>
          </b-navbar>
        </div>
        <div class="card-block">
          <div class="row">
            <div class="col-md-6">
              <h4>Ticket Types</h4>
              <div v-for="type in project.ticketTypes">
                <router-link :to="typeQuery(type)">
                  {{ type }}
                </router-link>
              </div>
            </div>
            <div class="col-md-6">
              <h4>Recent Activity</h4>
              <h1>TODO</h1>
            </div>
          </div>
        </div>
        <div class="card-block">
          <div class="row">
            <div class="col-3">
            </div>
            <div class="col-9" >
              <ticket-show :ticket
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <loading-spinner></loading-spinner>
    </div>
  </div>
</template>

<script>
 import Axios from 'axios'
 import LoadingSpinner from '@/components/General/LoadingSpinner'

 export default {
   name: 'project',

   components: {
     LoadingSpinner
   },

   computed: {
     ticketLink: function () {
       return '/queries?q=project = "' + this.project.key + '"'
     }
   },

   methods: {
     typeQuery: function (type) {
       return '/queries?q=project = "' + this.project.key + '" AND type = "' + type + '"'
     }
   },

   data: function () {
     return {
       project: {}
     }
   },

   created: function () {
     let url = '/api/projects/' + this.$route.params.key
     let inst = this

     Axios.get(url)
          .then((res) => {
            inst.project = res.data
          })
          .catch((err) => {
            if (err.response.status === 404) {
              this.$router.push('/404')
            }

            console.log('ERROR', err.response.data)
          })
   }
 }
</script>


<style lang="scss">
 .projectNav {
   background-color: #ddd;
 }
</style>
