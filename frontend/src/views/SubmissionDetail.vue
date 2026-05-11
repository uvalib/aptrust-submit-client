<template>
   <div class="details">
      <WaitSpinner v-if="submission.working" message="Please wait<br/>Loading submission details..." :overlay="true" />
      <template v-else>
         <h1>AP Trust Submission Detail</h1>
         <div class="info">
            <dl>
               <dt>Submission</dt>
               <dd>{{ submission.detail.collectionName }}</dd>
               <dt>Identifier</dt>
               <dd>{{ submission.detail.identifier }}</dd>
               <dt>Client</dt>
               <dd>{{ submission.detail.client.name }}</dd>
               <template v-if="submission.detail.client.approvalEmail">
                  <dt>Approval Email</dt>
                  <dd>{{ submission.detail.client.approvalEmail }}</dd>
               </template>
               <dt>Default Storage</dt>
               <dd>{{ submission.detail.storage }}</dd>

               <dt>Created</dt>
               <dd>{{ $formatDateTime(submission.detail.createdAt) }}</dd>

               <template v-if="submission.detail.approval">
                  <dt>Approved By</dt>
                  <dd>{{ submission.detail.approval.who }}</dd>
                  <dt>Approved On</dt>
                  <dd>{{ $formatDateTime(submission.detail.approval.createdAt) }}</dd>
               </template>
               <dt>Status</dt>
               <dd>
                  <div class="button-data">
                     <span class="caps">{{ submission.currentStatus.replace("-", " ") }}</span>
                     <StatusDialog />
                  </div>
               </dd>

               <dt>Failures</dt>
               <dd>
                  <div class="button-data">
                     <span v-if="submission.detail.failures.length == 0" class="none">None</span>
                     <template v-else>
                        <span class="error">{{ submission.detail.failures.length }}</span>
                        <FailuresDialog />
                     </template>
                  </div>
               </dd>

               <dt>Conflicts</dt>
               <dd>
                  <div class="button-data">
                     <span v-if="submission.detail.conflicts.length == 0" class="none">None</span>
                     <template v-else>
                        <span class="error">{{ submission.detail.conflicts.length }}</span>
                        <ConflictsDialog />
                     </template>
                  </div>
               </dd>

               <dt>Contents</dt>
               <dd>
                  <div class="button-data">
                     <div class="bag-summary" >
                        <div>{{ submission.detail.bagCount }} bag(s) containing {{ submission.detail.fileCount }} file(s).</div>
                        <div>Total size: {{ $formatFileSize(submission.detail.totalFileSize) }}</div>
                     </div>
                  </div>
               </dd>
            </dl>
            <ApprovePanel v-if="canApprove"/>
         </div>
         <BagsPanel />
      </template>
   </div>
</template>

<script setup>
import { onBeforeMount, ref, computed } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useRoute } from 'vue-router'
import WaitSpinner from '@/components/WaitSpinner.vue'
import StatusDialog from '@/components/StatusDialog.vue'
import ConflictsDialog from '@/components/ConflictsDialog.vue'
import FailuresDialog from '@/components/FailuresDialog.vue'
import ToggleSwitch from 'primevue/toggleswitch'
import BagsPanel from '@/components/BagsPanel.vue'
import ApprovePanel from '@/components/ApprovePanel.vue'

const submission = useSubmissionsStore()
const system = useSystemStore()
const route = useRoute()
const showBags = ref(false)

onBeforeMount( () => {
   submission.getSubmissionDetail( route.params.id )
   submission.getBags( route.params.id )  
})

const canApprove = computed( ()=> {
   if (submission.currentStatus != "pending-approval") return false 
   return system.canApproveSubmissions
})
</script>

<style lang="scss">
.details {
   margin: 0 auto 50px;
   min-height: 600px;
   text-align: left;

   .info {
      margin-bottom: 20px;
      display: flex;
      flex-flow: row wrap;
      gap: 20px;
      justify-content: flex-start;
      align-items: flex-start;
   }

   .p-panel {
      flex: 35%;
   }
   dl {
      grid-template-columns: max-content 2fr;
      display: inline-grid;
      grid-column-gap:  2rem;
      margin: 0;
      flex: 60%;

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
         .bag-summary {
            display: flex;
            flex-direction: column;
            gap: 5px;
         }
         .button-data {
            display: flex;
            flex-flow: row nowrap;
            align-items: center;
            justify-content: flex-start;
            gap: 30px;
            .labeled-toggle {
               display: flex;
               flex-flow: row nowrap;
               align-items: center;
               justify-content: flex-start;
               gap: 10px;   
            }
            span.caps {
               text-transform: capitalize
            }
            .p-button {
               padding: 2px 1rem;
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
      .info {
         flex-direction: column-reverse;
      }
   }
}
</style>
   