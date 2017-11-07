<template>
  <div class="container">
    <div v-if="user.username">
      <div class="row">
        <div class="col-md-3">
          <img v-bind:src="user.profilePicture" />
          <h3>{{ user.username }}</h3>
          <div>{{ user.fullName }}</div>
          <div>{{ user.email }}</div>
          <div><p>{{ user.bio }}</p></div>
        </div>
        <div class="col-md-9">
          <h1>Recent Activity</h1>
          <div>
            <template v-if="notifications.length > 0">
              <template v-for="notification in notifications">
                <notification :notification="notification" />
              </template>
            </template>
            <template v-else>
              <h4>This user doesn't have any activity yet!</h4>
            </template>
          </div>
          <h1>Roles</h1>
          <div v-for="project in projects" class="card role-card">
            <h5 class="card-block">
              {{ project.role }} of
              <router-link v-bind:to="'/projects/' + project.project">
                {{ project.project }}
              </router-link>
            </h5>
          </div>
          <div v-for="role in user.roles" class="card role-card">
            <h5 class="card-block">
              {{ role.role }} of
              <router-link v-bind:to="'/projects/' + role.project">
                {{ role.project }}
              </router-link>
            </h5>
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
 import LoadingSpinner from '@/components/General/LoadingSpinner'
 import Notification from '@/components/General/Notification'
 import Axios from 'axios'

 export default {
   name: 'user',
   components: {
     Notification,
     LoadingSpinner
   },

   watch: {
     leadOf: function () {
       this.projects = this.leadOf
                           .map((p) => { return {role: 'Lead', project: p.key} })
                           .concat(this.user.roles)
     }
   },

   data: function () {
     return {
       projects: [],
       leadOf: [],
       notifications: [],
       user: { roles: [] }
     }
   },

   created: function () {
     let inst = this
     let username = this.$route.params.username

     Axios.get('/api/users/' + username)
                        .then((res) => {
                          inst.user = res.data
                        })
                        .catch((err) => {
                          if (err.response.status === 404) {
                            this.$router.push('/404')
                          }

                          console.log('ERROR', err.response.data)
                        })

     Axios.get('/api/users/' + username + '/leadof')
                        .then((res) => {
                          inst.leadOf = res.data
                        })
                        .catch((err) => {
                          console.log('ERROR', err.response.data)
                        })

     Axios.get('/api/users/' + username + '/activity')
                        .then((res) => {
                          inst.notifications = res.data
                        })
                        .catch((err) => {
                          console.log('ERROR', err.response.data)
                        })
   }
 }
</script>

<style>
 .role-card .card-block {
   padding-top: 0.5rem;
 }
</style>
