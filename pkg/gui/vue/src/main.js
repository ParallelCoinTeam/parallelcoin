import Vue from 'vue'
import Kernel from './Kernel.vue'

Vue.config.productionTip = false

Vue.config.devtools = true;
// Vue.use(VueFormGenerator);

// Vue.prototype.$eventHub = new Vue(); 
 


new Vue({
  render: h => h(Kernel),
}).$mount('#kernel')
