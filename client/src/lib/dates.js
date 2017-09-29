import moment from 'moment'

export default {
  dateFormat (date) {
    return moment(date).format('YYYY-MM-DD hh:mm A')
  }
}
