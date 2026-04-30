<template>
   <div class="details">
      <WaitSpinner v-if="submission.working" message="Please wait<br/>Loading submission details..." :overlay="true" />
      <template v-else>
         <h1>Submission {{ submission.detail.collectionName }}</h1>  
         <div class="info">
            <dl>
               <dt>Client</dt>
               <dd>{{ submission.detail.client.name }}</dd>

               <dt>Created</dt>
               <dd>{{ $formatDateTime(submission.detail.createdAt) }}</dd>

               <dt>Storage</dt>
               <dd>{{ submission.detail.storage }}</dd>

               <dt>Status</dt>
               <dd>
                  <div class="status">
                     <span>{{ submission.currentStatus }}</span>
                     <StatusDialog :status="submission.detail.status" />
                  </div>
               </dd>

               <dt>Failures</dt>
               <dd>
                  <div class="errors">
                     <span v-if="submission.detail.failures.length == 0" class="none">None</span>
                     <template v-else>
                        <span class="error">{{ submission.detail.failures.length }}</span>
                     </template>
                  </div>
               </dd>

               <dt>Conflicts</dt>
               <dd>
                  <div class="errors">
                     <span v-if="submission.detail.conflicts.length == 0" class="none">None</span>
                     <template v-else>
                        <span class="error">{{ submission.detail.conflicts.length }}</span>
                        <ConflictsDialog :conflicts="submission.detail.conflicts" />
                     </template>
                  </div>
               </dd>

               <dt>Contents</dt>
               <dd>{{ submission.detail.bagCount }} bag(s) containing {{ submission.detail.fileCount }} file(s). Total size: {{ submission.totalSize }}</dd>
            </dl>
         </div>
      </template>
   </div>
</template>

<script setup>
import { onBeforeMount } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useRoute } from 'vue-router'
import WaitSpinner from '@/components/WaitSpinner.vue'
import StatusDialog from '@/components/StatusDialog.vue'
import ConflictsDialog from '@/components/ConflictsDialog.vue'

const submission = useSubmissionsStore()
const system = useSystemStore()
const route = useRoute()

onBeforeMount( () => {
   submission.getSubmissionDetail( route.params.id )
})
</script>

<style lang="scss">
.details {
   margin: 0 auto 50px;
   min-height: 600px;
   text-align: left;

   dl {
      grid-template-columns: max-content 2fr;
      display: inline-grid;
      grid-column-gap:  2rem;
      width: 100%;
      margin: 0;

      dt {
         font-weight: bold;
         text-align: right;
         padding: 0.5rem 0;
         white-space: nowrap;
      }
      dd {
         margin: 0;
         width: 100%;
         text-align: left;
         padding: 0.5rem 0;
         .none {
            font-style: italic;
            color: $uva-grey-A;
         }
         .status, .errors {
            display: flex;
            flex-flow: row nowrap;
            align-items: baseline;
            justify-content: flex-start;
            gap: 10px;
            span {
               text-transform: capitalize
            }
         }
         .error {
            font-weight: bold;
            color: $uva-red-A;
         }
      }
   }
}
@media only screen and (min-width: 768px) {
   .details {
       width: 90%;
   }
}
@media only screen and (max-width: 768px) {
   .details {
      width: 95%;
   }
}
</style>
   