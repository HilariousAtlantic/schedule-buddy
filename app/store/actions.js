import axios from 'axios';

export default {

  fetchTerms({commit}) {

    commit('REQUEST_TERMS');

    axios.get('/api/terms').then(response => {

      let terms = response.data;

      commit('RECEIVE_TERMS', terms);

    });

  },

  selectTerm({commit, state}, term) {

    commit('SELECT_TERM', term);

    let courses = state.coursesCache[term.id];

    if (courses) {

      commit('RECEIVE_CACHED_COURSES', courses);

    } else {

      commit('REQUEST_COURSES');

      axios.get('/api/courses?term='+term.id).then(response => {

        let courses = response.data.sort(

          (a, b) => (a.subject+a.number).localeCompare(b.subject+b.number)

        );

        commit('RECEIVE_COURSES', {term, courses});

      });

    }

  },

  setFilter({commit}, filter) {

    commit('SET_FILTER', filter);

  },

  selectCourse({commit, state}, course) {

    let detailedCourse = state.detailedCourseCache[course.id];

    if (!detailedCourse) {

      commit('REQUEST_DETAILED_COURSE');

      axios.get('/api/courses/'+course.id).then(response => {

        let detailedCourse = response.data;

        commit('RECEIVE_DETAILED_COURSE', detailedCourse);

        commit('SELECT_COURSE', detailedCourse);

      });

    } else {

      commit('SELECT_COURSE', course);

    }

  },

  deselectCourse({commit}, course) {

    commit('DESELECT_COURSE', course);

  },

  generateSchedules({commit, state}) {

    let courses = state.selectedCourses.join(',');

    let schedules = state.schedulesCache[courses];

    if (!schedules) {

      commit('REQUEST_SCHEDULES');

      axios.get('/api/schedules?courses='+courses).then(response => {

          let schedules = response.data;

          commit('RECEIVE_SCHEDULES', {courses, schedules});

      });

    } else {

      commit('RECEIVE_CACHED_SCHEDULES', schedules);

    }

  },

  increaseSchedulesIndex({commit}) {

    commit('INCREASE_SCHEDULES_INDEX');

  },

  decreaseSchedulesIndex({commit}) {

    commit('DECREASE_SCHEDULES_INDEX');

  }

}
