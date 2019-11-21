import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';

Vue.use(Router);

function lazyLoad(view) {
  return () => import(`@/views/${view}.vue`);
}

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'root',
      component: Home,
    },
    {
      path: '/about',
      name: 'about',
      component: () => lazyLoad('./views/About.vue'),
    },
    {
      path: '/contact',
      name: 'contact',
      component: () => lazyLoad('./views/Contact.vue'),
    },
    {
      path: '/events',
      name: 'events',
      component: () => lazyLoad('./views/Events.vue'),
    },
    {
      path: '/media',
      name: 'media',
      component: () => lazyLoad('./views/Media.vue'),
    },
    {
      path: '/members',
      name: 'members',
      component: () => lazyLoad('./views/Members.vue'),
    },
    {
      path: '/merch',
      name: 'merch',
      component: () => lazyLoad('./views/Merch.vue'),
    },
    {
      path: '/projects',
      name: 'projects',
      component: () => lazyLoad('./views/Projects.vue'),
    },
    {
      path: '/resources',
      name: 'resources',
      component: () => lazyLoad('./views/Resources.vue'),
    },
  ],
});
