<template>
  <div id="ticket-list-root">
    <div id="list-wrapper" v-if="tickets">
      <div>
        <h2>Columns</h2>
        <div v-for="column in columns">
          <span>{{ column.displayName ? column.displayName : humanizeColumnName(column.name) }}</span>
          <input type="checkbox" v-model="column.active" />
        </div>
      </div>
      <table id="ticket-list" class="table">
        <thead>
          <tr>
            <var v-for="column in columns">
              <th v-if="column.active">
                {{ column.displayName ? column.displayName : humanizeColumnName(column.name) }}
              </th>
            </var>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ticket in tickets">
            <var v-for="column in columns">
              <var v-if="column.active">
                <td v-if="ticket[column.name]">
                  {{ ticket[column.name] }}
                </td>
                <td v-else>
                  {{ getFieldValue(ticket, column.name) }}
                </td>
              </var>
            </var>
          </tr>
        </tbody>
      </table>
    </div>
    <h1 v-else>No tickets!</h1>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'ticket-list',
  methods: {
    getFieldValue: function (ticket, fieldName) {
      let field = ticket.fields.filter(f => f.name === fieldName)
      return field ? field.value : 'None'
    },

    humanizeColumnName: function (columnName) {
      return columnName
        // insert a space before all caps
        .replace(/([A-Z])/g, ' $1')
        // uppercase the first character
        .replace(/^./, function (str) { return str.toUpperCase() })
    }
  },

  watch: {
    tickets: function () {
      this.columns.concat(this.tickets[0].fields
        .map(f => { return { name: f.name, active: true } }))
    }
  },

  data: function () {
    return {
      columns: [
        {
          name: 'key',
          active: true
        },
        {
          name: 'summary',
          active: true
        },
        {
          name: 'createdDate',
          active: false
        },
        {
          name: 'updatedDate',
          active: false
        },
        {
          name: 'status',
          active: true
        },
        {
          name: 'project',
          active: true
        },
        {
          name: 'description',
          active: false
        },
        {
          name: 'assignee',
          active: false
        },
        {
          name: 'reporter',
          active: false
        },
        {
          name: 'labels',
          active: false
        },
        {
          name: 'ticketType',
          active: true
        }
      ]
    }
  },

  computed: mapState({
    tickets: function (state) {
      return state.tickets
    }
  })
}
</script>
