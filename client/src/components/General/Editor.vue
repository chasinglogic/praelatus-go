<template>
  <div class="md-editor card">
    <div class="card-block">
      <b-navbar>
        <b-nav is-nav-bar>
          <b-button @click="togglePreview">
            Preview
          </b-button>
        </b-nav>
      </b-navbar>
      <textarea v-model="text" v-bind:class="{ preview: preview }">
      </textarea>
      <div v-show="preview" class="preview card">
        <div class="card-block" v-html="renderedPreview" >
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
     }
   },

   props: {
     text: ''
   },

   data: function () {
     return {
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
</style>
