<template>
   <div class="dashboard">
      <h1>Browse AP Trust Submissions</h1>   
      <div class="content">
         <DataTable :value="submissionStore.searchHits" ref="hitstable" dataKey="id"
            selectionMode="single" v-model:selection="selectedSubmission" @update:selection="submissionClicked"
            v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
            stripedRows showGridlines responsiveLayout="scroll"
            :lazy="true" :paginator="true" :alwaysShowPaginator="true"
            @page="onPage($event)"  paginatorPosition="both" :rowsPerPageOptions="[30,50,100]" 
            :first="submissionStore.offset" :rows="submissionStore.pageSize" :totalRecords="submissionStore.total"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
            sortField="createdAt" :sortOrder="-1" @sort="onSort($event)" 
            :loading="submissionStore.working"
         >
            <template #header>
               <div class="search-controls">
                  <div class="search">
                     <Checkbox v-model="submissionStore.includeAutoApproved" inputId="autoapproved" binary @update:model-value="includeAutoApproveChanged()"/>
                     <label for="autoapproved">Include auto-approved submissions</label>
                  </div>
                  <div class="search">
                     <IconField iconPosition="left" class="query">
                        <InputIcon class="pi pi-search" />
                        <InputText v-model="submissionStore.query" @keypress="searchKeyPressed($event)" fluid aria-label="search submissions"/>
                     </IconField>
                     <Button severity="secondary" label="Reset" @click="resetSearchClicked()"/>
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
            <Column field="statusCreatedAt" header="Updated" class="nowrap" sortable>
               <template #body="slotProps">
                  <span>{{ $formatDateTime(slotProps.data.statusCreatedAt)  }}</span>
               </template>
            </Column>
         </DataTable>
      </div>
   </div>
</template>

<script setup>
import { onMounted, ref, computed, onBeforeMount } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useRouter, useRoute } from 'vue-router'
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
const route = useRoute()

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

onBeforeMount( () => {
   // pull query params and update store filter, query and pagination settings
   // they will be used in onMounted when submissions are requested
   if ( route.query.autoapproved ) {
      submissionStore.includeAutoApproved = true 
   }
   if ( route.query.q ) {
      submissionStore.query = route.query.q 
   }
   if ( route.query.offset ) {
      submissionStore.offset = parseInt(route.query.offset, 10)
   }
   if ( route.query.limit ) {
      submissionStore.pageSize = parseInt(route.query.limit, 10)
   }
   if ( route.query.sort ) {
      submissionStore.sortField = route.query.sort
   }
   if ( route.query.order ) {
      submissionStore.sortOrder = route.query.order
   }
   if ( route.query.filter ) { 
      let filterJson = JSON.parse(route.query.filter)
      filterJson.forEach( fj => {
         let bits = fj.split("=")
         submissionStore.filters.push( {field: bits[0], match: "equals", value: bits[1]})
      })
   }
})

onMounted( () => {
   submissionStore.getSubmissions()
})

const resetSearchClicked = (() => {
   let query = Object.assign({}, route.query)
   delete query.autoapproved
   delete query.q
   delete query.offset
   delete query.pageSize
   delete query.sortField
   delete query.sortOrder
   delete query.filter
   router.push({query})

   // clear the primevue datatables filter model
   for (let key in filters.value) {
      filters.value[key].value = null
   }
   
   submissionStore.resetSearch()
})

const searchToQueryParams = (() => {
   let query = Object.assign({}, route.query)
   delete query.autoapproved
   delete query.q
   delete query.offset
   delete query.pageSize
   delete query.sortField
   delete query.sortOrder
   delete query.filter
   if ( submissionStore.includeAutoApproved ) {
      query.autoapproved = 1
   }
   if ( submissionStore.query != "" ) {
      query.q = submissionStore.query
   }
   if (submissionStore.offset > 0 ) {
      query.offset = submissionStore.offset
   }
   query.limit = submissionStore.pageSize
   query.sort = submissionStore.sortField
   query.order = submissionStore.sortOrder
   if (submissionStore.filters.length > 0) {
      query.filter = submissionStore.filtersAsQueryParam
   }
   router.push({query})
   submissionStore.getSubmissions()
})

const includeAutoApproveChanged = (() => {
   searchToQueryParams()
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
   searchToQueryParams()
})

const onFilter = ((event) => {
   submissionStore.offset = 0
   submissionStore.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         submissionStore.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   searchToQueryParams()
})

const onPage = ((event) => {
   submissionStore.offset = event.first
   submissionStore.pageSize = event.rows
   searchToQueryParams()
})

const searchKeyPressed = ((event) => {
   submissionStore.offset = 0
   if (event.keyCode == 13) {
      searchToQueryParams()
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