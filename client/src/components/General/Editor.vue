<template>
  <div class="md-editor card">
    <div class="card-block">
      <div class="toolbar-wrapper">
        <div class="toolbar toolbar-left">
          <b-button v-b-tooltip.hover.auto title="Toggle Heading 1 Current Line"
            @click="toggleHeading(1)">
            <strong>
              H1
            </strong>
          </b-button>
          <b-button v-b-tooltip.hover.auto title="Toggle Heading 2 Current Line"
            @click="toggleHeading(2)">
            <strong>
              H2
            </strong>
          </b-button>
          <b-button v-b-tooltip.hover.auto title="Toggle Heading 3 Current Line"
            @click="toggleHeading(3)">
            <strong>
              H3
            </strong>
          </b-button>
          <b-button v-b-tooltip.hover.auto title="Toggle Bold on Selection" @click="toggleBold">
            <strong>
              B
            </strong>
          </b-button>
          <b-button v-b-tooltip.hover.auto title="Toggle Italic on Selection"
            @click="toggleItalic">
            <i class="fa fa-italic"></i>
          </b-button>
          <b-button v-b-tooltip.hover.auto title="Toggle Code on Selection" @click="toggleCode">
            <i class="fa fa-code"></i>
          </b-button>
          <b-button v-b-tooltip.hover.auto title="Start List" @click="toggleList(false)">
            <i class="fa fa-list-ul"></i>
          </b-button>
          <b-button v-b-tooltip.hover.auto title="Start Numbered List" @click="toggleList(true)">
            <i class="fa fa-list-ol"></i>
          </b-button>
        </div>
        <div class="toolbar toolbar-right">
          <b-button v-b-tooltip.hover.auto title="Preview Markdown" @click="togglePreview">
            <i class="fa fa-eye"></i>
          </b-button>
          <b-button v-b-tooltip.hover.auto title="Help" @click="toggleHelp">
            <i class="fa fa-question-circle"></i>
          </b-button>
        </div>
      </div>
      <div class="toolbar-wrapper">
        <textarea style="resize: none"
          v-model="text"
          @keydown="handleKeyPress"
          v-bind:id="id"
          v-bind:class="{ preview: preview }">
        </textarea>
        <div v-show="preview" class="preview card">
          <div class="card-block" v-html="renderedPreview" >
          </div>
        </div>
        <div v-show="showHelp" class="preview card">
          <div class="card-block">
            <h5 class="text-center">Editor Help</h5>
            <p>
              This editor accepts Markdown using the CommonMark standard.
              for more information on writing Markdown read this
              <a href="http://commonmark.org/help/">document</a>
            </p>
            <h6>Editor Keybindings</h6>
            <ul>
              <li>CTRL+h: Toggle this help.</li>
              <li>CTRL+p: Toggle rendered preview.</li>
              <li>CTRL+Enter: Submit for the current form.</li>
              <li>CTRL+{1, 2, 3, 4, 5}: Toggle a heading {1, 2, 3, 4, 5} on the current line.</li>
              <li>CTRL+i Toggle italics on the selected text.</li>
              <li>CTRL+b Toggle bold on the selected text.</li>
              <li>CTRL+y Start a multi line code block.</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
 import Markdown from '@/lib/markdown'

 export default {
   name: 'editor',

   computed: {
     renderedPreview: function () {
       return Markdown.render(this.text)
     }
   },

   watch: {
     value: function (newVal, oldVal) {
       if (newVal !== oldVal) {
         this.text = newVal
       }
     },

     text: function (newVal, oldVal) {
       if (newVal !== oldVal) {
         this.$emit('input', newVal)
       }
     }
   },

   methods: {
     handleKeyPress: function (ev) {
       const boundKey = this.boundKeys[ev.key]
       // Prevent the event from being passed up the chain
       ev.stopPropagation()
       if (boundKey) {
         let txta = document.getElementById(this.id)
         let curpos = txta.selectionStart
         boundKey(ev, curpos)
       }
     },

     togglePreview: function () {
       this.preview = this.preview === false
       document.getElementById(this.id).focus()
     },

     toggleHelp: function () {
       this.showHelp = this.showHelp === false
       document.getElementById(this.id).focus()
     },

     findPreviousNewline: function (i) {
       if (i === null) {
         return null
       }

       let newPos = this.text[i] === '\n' ? i - 1 : i
       while (newPos !== 0 && this.text[newPos] !== '\n') {
         newPos = newPos - 1
       }
       return newPos
     },

     toggleHeading: function (headingSize) {
       let txta = document.getElementById(this.id)
       let curpos = txta.selectionStart
       let prevNewline = this.findPreviousNewline(curpos)
       // When at the beginning of the document (index 0) just use that
       // otherwise extend the "beginning half" to include the previous
       // newline character
       let x = prevNewline === 0 ? prevNewline : prevNewline + 1

       // This if condition checks if a heading of the given size already
       // exists on this line, if so we "toggle" it off.
       if (this.text[x] === '#' &&
           this.text[x + (headingSize - 1)] === '#') {
         let beg = this.text.slice(0, x)
         // Cut off the number of heading characters plus 1 for the space
         let end = this.text.slice(x + headingSize + 1, this.text.length)
         // Join our now cut string
         this.text = beg + end
         // Else if we find a newline before a heading then insert a new
         // heading of the given size
       } else {
         // Add a space for proper markdown formatting
         let heading = '#'.repeat(headingSize) + ' '
         this.insertAt(x, heading)
       }

       txta.focus()
     },

     toggleList: function (ordered) {
       let txta = document.getElementById(this.id)
       let curpos = txta.selectionStart
       let prevNewline = this.findPreviousNewline(curpos)
       // When at the beginning of the document (index 0) just use that
       // otherwise extend the "beginning half" to include the previous
       // newline character
       let x = prevNewline === 0 ? prevNewline : prevNewline + 1
       if (ordered) {
         this.insertAt(x, ' 1. ')
       } else {
         this.insertAt(x, ' - ')
       }

       txta.focus()
     },

     wrapText: function (start, end, char) {
       // Grab everything before the selection
       this.text = this.text.slice(0, start) +
                   // Grab the actual selection and wrap it in *
                   char + this.text.slice(start, end) + char +
                   // Grab the last of the text and concat each section to form
                   // new text
                   this.text.slice(end, this.text.length)
     },

     unwrapText: function (start, end, char) {
       let newStart = start - char.length
       let newEnd = end + char.length + 1
       let newSelectStart = start
       let newSelectEnd = end + 1

       if (this.text[start] === char[0]) {
         newStart = newStart - 1
         newSelectStart = newSelectStart + 1
       }

       if (this.text[end] === char[0]) {
         newEnd = newEnd - 1
         newSelectEnd = newSelectEnd - 1
       }

       let beg = this.text.slice(0, newStart)
       let mid = this.text.slice(newSelectStart, newSelectEnd)
       let rest = this.text.slice(newEnd, this.text.length)

       this.text = beg + mid + rest
     },

     isSurroundedBy: function (start, end, char) {
       return (this.text.slice(start - char.length, start) === char &&
               this.text.slice(end, end + char.length) === char)
     },

     insertAt: function (x, text) {
       // Grab the front of our string
       let beg = this.text.slice(0, x) + text
       // Grab everything after where we inserted.
       let end = this.text.slice(x, this.text.length)
       // Set the editor text equal to our newly eddited halves
       this.text = beg + end
     },

     toggleItalic: function () {
       let txta = document.getElementById(this.id)
       let start = txta.selectionStart
       let end = txta.selectionEnd
       // Check if text is already wrapped in asterisks
       if (this.isSurroundedBy(start, end, '*')) {
         this.unwrapText(start, end, '*')
       } else {
         this.wrapText(start, end, '*')
       }

       txta.focus()
     },

     toggleBold: function () {
       let txta = document.getElementById(this.id)
       let start = txta.selectionStart
       let end = txta.selectionEnd
       // Check if text is already wrapped in double asterisks
       if (this.isSurroundedBy(start, end, '**')) {
         this.unwrapText(start, end, '**')
       } else {
         this.wrapText(start, end, '**')
       }

       txta.focus()
     },

     toggleCode: function () {
       let txta = document.getElementById(this.id)
       let start = txta.selectionStart
       let end = txta.selectionEnd

       if (this.isSurroundedBy(start, end, '```')) {
         this.unwrapText(start, end, '```')
       } else if (this.isSurroundedBy(start, end, '`')) {
         this.unwrapText(start, end, '`')
       } else if (this.text[start] === '\n' &&
                  (this.text[end] === '\n' || this.text[end] === undefined)) {
         this.wrapText(start, end, '```')
       } else {
         this.wrapText(start, end, '`')
       }

       txta.focus()
     }
   },

   props: [
     'value'
   ],

   data: function () {
     return {
       // Generate a unique ID for the text area in the event that the editor
       // component is on the same page multiple times.
       id: Math.random().toString(36).substring(2, 5),
       preview: false,
       showHelp: false,
       text: this.value,
       // Use arrow functions in these definitions so that "this"
       // is properly set.
       boundKeys: {
         '1': (ev, curpos) => {
           if (!ev.ctrlKey) {
             return
           }

           this.toggleHeading(1)
         },

         '2': (ev, curpos) => {
           if (!ev.ctrlKey) {
             return
           }

           this.toggleHeading(2)
         },

         '3': (ev, curpos) => {
           if (!ev.ctrlKey) {
             return
           }

           this.toggleHeading(3)
         },

         '4': (ev, curpos) => {
           if (!ev.ctrlKey) {
             return
           }

           this.toggleHeading(4)
         },

         '5': (ev, curpos) => {
           if (!ev.ctrlKey) {
             return
           }

           this.toggleHeading(5)
         },

         Enter: (ev, curpos) => {
           if (ev.ctrlKey) {
             ev.preventDefault()
             this.$emit('submit')
             return
           }

           let prevNewline = this.findPreviousNewline(curpos)
           prevNewline = prevNewline === 0 ? prevNewline : prevNewline + 1
           let slice = this.text.slice(prevNewline, curpos)

           let listRgx = /^ [-*] (.*?)/
           let numListRgx = /^( *)([0-9]*)\. (.*?)/

           if (listRgx.exec(slice)) {
             ev.preventDefault()
             let m = listRgx.exec(slice)
             this.insertAt(curpos, '\n' + m[0])
           } else if (numListRgx.exec(slice)) {
             ev.preventDefault()
             let m = numListRgx.exec(slice)
             let curNum = parseInt(m[2])
             this.insertAt(curpos, '\n' + m[1] + (curNum + 1).toString() + '. ')
           }
         },

         i: (ev, curpos) => {
           if (ev.ctrlKey) {
             ev.preventDefault()
             this.toggleItalic()
           }
         },

         y: (ev, curpos) => {
           if (ev.ctrlKey) {
             ev.preventDefault()
             this.insertAt(ev.target.selectionStart, '\n```\n')
           }
         },

         b: (ev, curpos) => {
           if (ev.ctrlKey) {
             ev.preventDefault()
             this.toggleBold()
           }
         },

         p: (ev, curpos) => {
           if (ev.ctrlKey) {
             ev.preventDefault()
             this.togglePreview()
           }
         },

         h: (ev, curpos) => {
           if (ev.ctrlKey) {
             ev.preventDefault()
             this.toggleHelp()
           }
         }
       }
     }
   }
 }
</script>

<style lang="scss">
 @import './src/assets/styles/globals.scss';

 $editor-height: 20rem;

 .md-editor .toolbar-wrapper textarea {
   width: 100%;
   height: $editor-height;
 }

 .md-editor .card {
   height: $editor-height;
   overflow-y: scroll;
 }

 .md-editor {
   width: 100%;
   margin-left: auto;
   margin-right: auto;
 }

 .md-editor {
   background-color: $faded-grey;
 }

 .md-editor .card-block {
   display: flex;
   flex-direction: column;
 }

 .toolbar {
   padding: 0.5em;
 }

 .toolbar-left {
   flex-grow: 2;
   text-align: left;
 }

 .toolbar-right {
   text-align: right;
 }

 .toolbar div {
   display: inline-block;
 }

 .toolbar-wrapper {
   display: flex;
   width: 100%;
 }

 .toolbar-wrapper .preview {
   min-width: 50%;
   text-align: left;
   padding: 0.25rem;
 }
</style>
