<template>
  <div id="loginForm" class="container">
    <h1 class="text-center" v-if="register">Register</h1>
    <h1 class="text-center" v-else>Login</h1>
    <b-form @submit="onSubmit">
      <b-form-group id="usernameInputGroup"
        label="Username" label-for="username">
        <b-form-input id="username"
          type="text" v-model="form.username" required
          placeholder="Username"
        ></b-form-input>
      </b-form-group>

      <b-form-group id="passwordInputGroup"
        label="Password" label-for="password">
        <b-form-input id="password"
          type="password" v-model="form.password" required
          placeholder="Password"
        ></b-form-input>
      </b-form-group>


      <template v-if="register">
        <b-form-group id="fullNameInputGroup"
          label="Full Name" label-for="fullName">

          <b-form-input id="full name"
            type="text" v-model="form.fullName" required
            placeholder="Full Name"
          ></b-form-input>
        </b-form-group>

        <b-form-group id="emailInputGroup"
          label="Your Email Address" label-for="email">

          <b-form-input id="email"
            type="text" v-model="form.email"
            placeholder="Email" required
          ></b-form-input>

        </b-form-group>
      </template>

      <b-button type="submit" variant="primary">Submit</b-button>
    </b-form>
  </div>
</template>

<script>
 import Axios from 'axios'

 export default {
   name: 'login',
   props: {
     'register': false
   },
   methods: {
     onSubmit: function (e) {
       e.preventDefault()

       let url = '/api/tokens'
       if (this.register) {
         url = '/api/users'
       }

       console.log('url', url)

       let inst = this
       Axios.post(url, inst.form)
            .then(function (resp) {
              localStorage['x-praelatus-token'] = resp.headers['x-praelatus-token']
              inst.$store.commit('login', {
                'token': resp.headers['x-praelatus-token'],
                'user': resp.data
              })

              inst.$router.push('/')
            })
            .catch(function (err) {
              if (err.response) {
                if (err.response.status === 404) {
                  inst.errors = ['No user with that username exists']
                } else if (err.response.status === 403) {
                  inst.errors = ['Invalid password.']
                } else {
                  inst.errors = [err.response.data.msg]
                }
              } else {
                console.log('Unhandled Error:', err)
              }
            })
     }
   },
   data: function () {
     return {
       'errors': [],
       'form': {
         'password': '',
         'username': '',
         'fullName': '',
         'email': ''
       }
     }
   }
 }
</script>

<style lang="scss">
 #loginForm form label {
   font-weight: bold;
 }

 #loginForm form button {
   width: 100%;
 }

 #loginForm form {
   text-align: left;
   max-width: 700px;
   margin-left: auto;
   margin-right: auto;
 }
</style>
