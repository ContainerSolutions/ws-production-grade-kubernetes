// const url = process.env.AVATAR_URL + '/store';
// const url = 'http://k8s-real-demo/store'
// const url = 'http://localhost:30697/webapp/store'
const url = '/webapp/store'
const iconUrl = '/webapp/avatar'

Vue.component('RecordList', {
  props: ['records'],

  template: `
    <ul>
      <li v-for="record in records"> 
        <record :record="record" readonly></record>
      </li>
      <li>
        <record :record='{}'></record>
      </li>
    </ul>
  `
});

Vue.component('Record', {
  props: ['record', 'readonly'],

  template: `
    <div>
      <input v-model="record.Key" type="text" name="key" v-bind:readonly="readonly">
      <input v-model="record.Value" type="text" name="key">
      <button v-on:click="updateRecord(record)">Create/Update</button>
      <button v-on:click="deleteRecord(record)">Delete</button>
    </div>
  `,

  methods: {
    deleteRecord: function(record) {
      return axios.delete(url, {data: {Key: record.Key}})
        .then((response) => {
          vm.getRecords()
        }
      );
    },
    updateRecord: function(record) {
      return axios.put(url, {Key: record.Key, Value: record.Value})
        .then((response) => {
          vm.getRecords()
        }
      );
    },
    generateIcon: function() {
      return axios.get(iconUrl)
        .then((response) => {
          // Display Icon?
        })
    }
  }
});

const vm = new Vue({
  el: '#app',
  data: {
    records: []
  },

  methods: {
    getRecords: function() {
      return axios.get(url)
        .then((response) => {
          this.records = response.data
      });
    },
  },

  mounted: function() {
    this.getRecords();
  }

});
