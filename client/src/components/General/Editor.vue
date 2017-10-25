<template>
  <div class="md-editor card">
    <div class="card-block">
      <div class="toolbar-wrapper">
        <div class="toolbar toolbar-left">
          <b-button @click="toggleHeading(1)">
            H1
          </b-button>
          <b-button @click="toggleHeading(2)">
            H2
          </b-button>
          <b-button @click="toggleBold">
            <strong>
              B
            </strong>
          </b-button>
          <b-button @click="toggleItalic">
            <em>
              I
            </em>
          </b-button>
          <b-button @click="toggleCode">
            <>
          </b-button>
        </div>
        <div class="toolbar toolbar-right">
          <b-button @click="togglePreview">
            Preview
          </b-button>
        </div>
      </div>
      <div class="toolbar-wrapper">
        <textarea style="resize: none"
          v-model="text"
          v-bind:id="id"
          v-bind:class="{ preview: preview }">
        </textarea>
        <div v-show="preview" class="preview card">
          <div class="card-block" v-html="renderedPreview" >
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

   methods: {
     togglePreview: function () {
       if (this.preview) {
         this.preview = false
       } else {
         this.preview = true
       }
     },

     findPreviousNewline: function (i) {
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

       // This if condition checks if a heading of the given size already
       // exists on this line, if so we "toggle" it off.
       if (this.text[prevNewline + 1] === '#' &&
           this.text[prevNewline + headingSize] === '#') {
         // When at the beginning of the document (index 0) just use that
         // otherwise extend the "beginning half" to include the previous
         // newline character
         let x = prevNewline === 0 ? prevNewline : prevNewline + 1
         let beg = this.text.slice(0, x)
         // Cut off the number of heading characters plus 1 for the space
         let end = this.text.slice(x + headingSize + 1, this.text.length)
         // Join our now cut string
         this.text = beg + end
         return
         // Else if we find a newline before a heading then insert a new
         // heading of the given size
       } else if (this.text[prevNewline] === '\n' || prevNewline === 0) {
         // If we are at the beginning of the doc use that, otherwise
         // increment by one so we include the newline at the beginning and
         // exlude if from the end.
         let x = prevNewline === 0 ? prevNewline : prevNewline + 1
         // Add a space for proper markdown formatting
         let heading = '#'.repeat(headingSize) + ' '
         this.insertAt(x, heading)
         return
       }
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
     },

     toggleCode: function () {
       let txta = document.getElementById(this.id)
       let start = txta.selectionStart
       let end = txta.selectionEnd

       console.log(this.text[start], this.text[start - 1])
       console.log(this.text[end], this.text[end + 1])

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
     }
   },

   props: {
     startingText: ''
   },

   mounted: function () {
     if (this.startingText && this.startingText !== '') {
       this.text = this.startingText
     }
   },

   data: function () {
     return {
       // Generate a unique ID for the text area in the event that the editor
       // component is on the same page multiple times.
       id: Math.random().toString(36).substring(2, 5),
       text: '',
       preview: false
     }
   }
 }
</script>

<style lang="scss">
 @import './src/assets/styles/globals.scss';

 $editor-height: 20rem;

 .md-editor textarea {
   width: 100%;
   height: $editor-height;
 }

 .md-editor .card {
   height: $editor-height;
   overflow-y: scroll;
 }

 .md-editor {
   max-width: 60rem;
 }

 .md-editor .preview {
   width: 48%;
   vertical-align: top;
   display: inline-block;
   text-align: left;
 }

 .md-editor {
   background-color: $faded-grey;
 }

 .toolbar {
   padding: 0.5em;
   display: inline-block;
 }

 .toolbar-left {
   text-align: left;
 }

 .toolbar-right {
   text-align: right;
 }

 .toolbar div {
   display: inline-block;
 }

 .toolbar-wrapper {
   display: block;
 }
</style>
