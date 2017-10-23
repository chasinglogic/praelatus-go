import Showdown from 'showdown'
const converter = new Showdown.Converter()

export default {
  render: function (text) {
    return converter.makeHtml(text)
  }
}
