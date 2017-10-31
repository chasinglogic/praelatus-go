<template>
  <div class="container">
    <h1>Create Ticket</h1>
    <div class="alert alert-danger" v-show="errorMsg !== null">
      {{ errorMsg }}
    </div>
    <b-form @submit="createTicket">
      <b-form-group id="projectSelection"
        label="Project:" label-for="projectSelector" >
        <b-form-select id="projectSelector"
          :options="projects" required
          v-model="ticket.project">
        </b-form-select>
      </b-form-group>

      <b-form-group id="ticketTypeSelection"
        label="Ticket Type:" label-for="ticketTypeSelector" >
        <b-form-select id="ticketTypeSelector"
          :options="ticketTypes" required
          v-model="ticket.type">
        </b-form-select>
      </b-form-group>

      <b-form-group id="summary"
        label="Summary:" label-for="summary" >
        <b-form-input id="summary" required
          v-model="ticket.summary">
        </b-form-input>
      </b-form-group>

      <b-form-group id="description"
        label="Description:" label-for="description" >
        <editor id="description"
          v-model="ticket.description">
        </editor>
      </b-form-group>

      <template v-for="field in ticket.fields">
        <b-form-group :id="field.name"
          :label="field.name + ':'" :label-for="field.name" >
          <template v-if="field.dataType === 'STRING'">
            <b-form-input :id="field.name"
              v-model="field.value">
            </b-form-input>
          </template>
          <template v-else-if="field.dataType === 'FLOAT'">
            <b-form-input :id="field.name"
              type="number"
              step="0.01"
              v-model="field.value">
            </b-form-input>
          </template>
          <template v-else-if="field.dataType === 'INT'">
            <b-form-input :id="field.name"
              type="number"
              v-model="field.value">
              </b-form-input>
          </template>
          <template v-else-if="field.dataType === 'DATE'">
            <b-form-input :id="field.name"
              type="datetime-local"
              v-model="field.value">
              </b-form-input>
          </template>
          <template v-else-if="field.dataType === 'OPT'">
            <b-form-select :id="field.name"
              :options="field.options"
              v-model="field.value">
            </b-form-select>
          </template>
          <template v-else>
            <div style="color: red; font-weight: bold">
              Error Rendering Custom Field Input
            </div>
          </template>
        </b-form-group>
      </template>

      <b-btn type="submit" class="form-control" variant="success">
        Create Ticket
      </b-btn>
    </b-form>
  </div>
</template>

<script>
 import Axios from 'axios'
 import Editor from '@/components/General/Editor'

 export default {
   name: 'ticket-create',

   components: {
     Editor
   },

   computed: {
     currentUser: function () {
       return this.$store.getters.currentUser
     },

     selectedProjectKey: function () {
       return this.ticket.project
     },

     selectedProject: function () {
       return this.projectData.filter(p => p.key === this.ticket.project)[0]
     },

     selectedType: function () {
       return this.ticket.type
     }
   },

   watch: {
     selectedProject: function () {
       this.getFieldSchemeForProject()
     },

     selectedType: function () {
       this.updateFields()
     }
   },

   methods: {
     createTicket: function () {
       let url = '/api/tickets'
       let inst = this

       Axios.post(url, this.ticket)
            .then((res) => {
              console.log(res.data)
              // inst.$router.push('/tickets/' + res.data.key)
            })
            .catch((err) => {
              inst.errorMsg = err.response.data.message
            })
     },

     updateFields: function () {
       if (this.fieldScheme === null ||
           this.fieldScheme === undefined ||
           this.fieldScheme.fields === null ||
           this.fieldScheme.fields === undefined) {
         return
       }

       if (this.fieldScheme.fields[this.ticket.type]) {
         this.ticket.fields = this.fieldScheme.fields[this.ticket.type]
         return
       }

       this.ticket.fields = this.fieldScheme.fields['']
     },

     getFieldSchemeForProject: function () {
       let url = '/api/fieldschemes/' + this.selectedProject.fieldScheme
       let inst = this

       Axios.get(url)
            .then((res) => {
              inst.fieldScheme = res.data
              this.updateFields()
            })
            .catch(console.log)
     },

     updateTicketTypes: function () {
       this.ticketTypes = this.selectedProject.ticketTypes
     }
   },

   data: function () {
     return {
       projects: [],
       errorMsg: null,
       projectData: [],
       ticketTypes: [],
       fieldScheme: {},
       ticket: {
         reporter: this.$store.getters.currentUser.username,
         type: '',
         summary: '',
         description: '',
         assignee: '',
         labels: [],
         fields: []
       }
     }
   },

   props: {
     project: {
       type: String,
       default: null
     },
     ticketType: {
       type: String,
       default: null
     }
   },

   created: function () {
     if (this.currentUser === undefined || this.currentUser === null) {
       this.$router.push({
         path: '/login',
         query: { to: this.$router.currentRoute.path }
       })
     }

     let url = '/api/projects?permission=CREATE_TICKET'
     let inst = this
     Axios.get(url)
          .then((res) => {
            inst.projectData = res.data
            inst.projects = res.data.map(p => p.key)
            if (inst.projectData.length === 0) { return }

            let preselectProject = this.$route.query.project
                    ? this.$route.query.project : this.project
            if (preselectProject) {
              let reqProject = inst.projects
                                   .filter(p => p === preselectProject)
              inst.ticket.project = reqProject.length > 0
                                  ? reqProject[0] : inst.projects[0]
            } else {
              inst.ticket.project = inst.projects[0]
            }

            this.updateTicketTypes()

            let preselectTicketType = this.$route.query.ticketType
                    ? this.$route.query.ticketType : this.ticketType
            if (preselectTicketType) {
              let reqType = inst.ticketTypes
                                .filter(t => t === preselectTicketType)
              inst.ticket.type = reqType.length > 0
                               ? reqType[0] : inst.ticketTypes[0]
            } else {
              inst.ticket.type = inst.ticketTypes[0]
            }

            this.getFieldSchemeForProject()
          })
          .catch((err) => {
            // TODO: Visually que the user that there's been an error
            console.log('ERROR', err)
          })
   }
 }
</script>

<style>
 form {
   text-align: left;
 }

 form label {
   font-weight: bold;
 }
</style>
