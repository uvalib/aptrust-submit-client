<template>
   <div class="dashboard">
      <h1>Browse AP Trust Submissions</h1>   
      <div class="content">
         <DataTable :value="submissionStore.searchHits" ref="hitstable" dataKey="id"
            selectionMode="single" v-model:selection="selectedSubmission" @update:selection="submissionClicked"
            v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
            stripedRows showGridlines responsiveLayout="scroll"
            :lazy="true" :paginator="true" :alwaysShowPaginator="true"
            @page="onPage($event)"  paginatorPosition="both"
            :first="submissionStore.offset" :rows="submissionStore.pageSize" :totalRecords="submissionStore.total"
            paginatorTemplate="PrevPageLink CurrentPageReport NextPageLink"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
            sortField="createdAt" :sortOrder="-1" @sort="onSort($event)" 
            :loading="submissionStore.working"
         >
            <template #header>
               <div class="search-controls">
                  <div class="search">
                     <Checkbox v-model="submissionStore.includeAutoApproved" inputId="autoapproved" binary @update:model-value="submissionStore.getSubmissions()"/>
                     <label for="autoapproved">Include auto-approved submissions</label>
                  </div>
                  <div class="search">
                     <IconField iconPosition="left" class="query">
                        <InputIcon class="pi pi-search" />
                        <InputText v-model="submissionStore.query" @keypress="searchKeyPressed($event)" fluid aria-label="search submissions"/>
                     </IconField>
                     <Button severity="secondary" label="Reset" @click="submissionStore.resetSearch()"/>
                  </div>
               </div>
            </template>
            <Column field="status" header="Status" filterField="status" :showFilterMatchModes="false" class="status-col">
               <template #filter="{ filterModel, filterCallback }">
                  <Select v-model="filterModel.value" @change="filterCallback()" :options="system.submissionStatuses" placeholder="Select a status" />
               </template>
               <template #body="slotProps">
                  <div :class="statusClass(slotProps.data.status)" v-if="slotProps.data.status">
                     <span>{{ slotProps.data.status.replace('-', ' ') }}</span>
                     <i v-if="slotProps.data.status=='pending-approval'" class="pi pi-spinner pi-spin" />
                     <i v-if="slotProps.data.status=='error'" class="pi pi-times-circle" />
                     <i v-if="slotProps.data.status=='abandoned'" class="pi pi-check-circle" />
                  </div>
                  <span v-else class="none">Unknown</span>   
               </template>
            </Column>
            <Column field="collectionName" header="Submission">
               <template #body="slotProps">
                  <span v-if="slotProps.data.collectionName">{{ slotProps.data.collectionName }}</span>
                  <span v-else class="none">Undefined</span> 
               </template>
            </Column>
            <Column field="client" header="Source" filterField="client" :showFilterMatchModes="false">
               <template #filter="{ filterModel, filterCallback }">
                  <Select v-model="filterModel.value" @change="filterCallback()" :options="clients" placeholder="Select a client" />
               </template>
            </Column>
            <Column field="createdAt" header="Created" class="nowrap" sortable>
               <template #body="slotProps">{{ $formatDateTime(slotProps.data.createdAt) }}</template>
            </Column>
            <Column field="statusCreatedAt" header="Uodated" class="nowrap" sortable>
               <template #body="slotProps">
                  <span>{{ $formatDateTime(slotProps.data.statusCreatedAt)  }}</span>
               </template>
            </Column>
         </DataTable>
      </div>
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useRouter } from 'vue-router'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import Checkbox from 'primevue/checkbox'
import Select from 'primevue/select'
import { FilterMatchMode } from '@primevue/core/api'

const submissionStore = useSubmissionsStore()
const system = useSystemStore()
const router = useRouter()

const selectedSubmission = ref()

const filters = ref({
   client: { value: null, matchMode: FilterMatchMode.EQUALS },
   status: { value: null, matchMode: FilterMatchMode.EQUALS },
   storage: { value: null, matchMode: FilterMatchMode.EQUALS },
});

const clients = computed( () => {
   let out = [] 
   system.clients.forEach( c=> {
      out.push(c.name)
   })
   return out
})

onMounted( () => {
   submissionStore.getSubmissions()
})

const statusClass = ( (status)=> {
   if (status == "pending-approval") return "status pending"
   if (status == "complete" || status == "abandoned") return "status done"
   if (status == "error") return "status error"
   return  "status"
})

const submissionClicked = (() => {
   router.push(`/submissions/${selectedSubmission.value.identifier}`)
})

const onSort = ((event) => {
   submissionStore.sortField = event.sortField
   if (event.sortOrder == 1) {
      submissionStore.sortOrder = "asc"

   } else if (event.sortOrder == -1) {
      submissionStore.sortOrder = "desc"
   }
   submissionStore.getSubmissions()
})

const onFilter = ((event) => {
   submissionStore.offset = 0
   submissionStore.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         submissionStore.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   submissionStore.getSubmissions()
})

const onPage = ((event) => {
   submissionStore.offset = event.first
   submissionStore.getSubmissions()
})

const searchKeyPressed = ((event) => {
   submissionStore.offset = 0
   if (event.keyCode == 13) {
      submissionStore.getSubmissions()
   }
})

</script>

<style lang="scss">
.dashboard {
   margin: 0 auto 50px;
   min-height: 600px;
   text-align: left;
   .none {
      color: $uva-grey-A;
      font-style: italic;
   }
   .nowrap {
      white-space: nowrap;
   }
   .status-col {
      width: 225px;
   }
   .status {
      text-transform: capitalize;
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      align-items: center;
      gap: 0.5rem;
      padding: 0.25rem .75rem; 
      background-color: $uva-grey-200;
      border: 1px solid $uva-grey;
      border-radius: 1rem;
   }
   .status.error {
       background-color: $uva-red-A;
       border-color: $uva-red-B;
       color: white;
       font-weight: bold;
   }
   .status.pending {
       background-color: $uva-yellow-100;
       border-color: $uva-yellow;
   }
   .status.done {
       background-color: $uva-green-A;
       border-color: $uva-green-B;
       color: white;
       font-weight: bold;
   }
   .search-controls {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      .search {
         display: flex;
         flex-flow: row nowrap;   
         align-items: center;
         justify-content: flex-start;
         gap: 10px;
      }
   }
}
@media only screen and (min-width: 768px) {
   .dashboard {
       width: 90%;
   }
}
@media only screen and (max-width: 768px) {
   .dashboard {
      width: 95%;
   }
}
</style>