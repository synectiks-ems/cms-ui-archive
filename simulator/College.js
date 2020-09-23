module.exports = {
  institutes: [
    { id: 0, code: 'JNU', name: 'Jawaharlal Nehru University', year: '1974'},
    {id: 1, code: 'RRK', name: 'Roorke Engineering College', year: '1976'},
  ],
  departments: [
    {id: 0, code: 'ECM', name: 'Electronics & Communication', description: 'desc'},
    {id: 1, code: 'COMP', name: 'Computer Science', description: 'description'},
  ],
  programs: [
    {id: 0, name: 'Bachlor of Engineering', shortname:'BE', description: 'desc'},
    {id: 1, name: 'Master of Engineering', shortname:'ME', description: 'description'},
  ],
  academiccalenders: [
    {id: 0, batchyear:2014, startdate:  new Date('2014-03-21T16:15:34.479Z'), enddate:new Date('2014-09-21T16:15:34.479Z'), program_id:0, semester: 'First'},
    {id: 1, batchyear:2014, startdate:  new Date('2014-09-21T16:15:34.479Z'), enddate:new Date('2015-03-21T16:15:34.479Z'), program_id:0, semester: 'Second'},
    {id: 2, batchyear:2015, startdate:  new Date('2015-03-21T16:15:34.479Z'), enddate:new Date('2015-09-21T16:15:34.479Z'), program_id:0, semester: 'First'},
    {id: 3, batchyear:2015, startdate:  new Date('2015-09-21T16:15:34.479Z'), enddate:new Date('2016-03-21T16:15:34.479Z'), program_id:0, semester: 'Second'},
  ],
  offers: [
    {id: 0, institute_id: 0, program_id:0, department_id:0, active:1, desc: 'Institute-Programme-Department mapping'},
    {id: 1, institute_id: 0, program_id:0, department_id:1, active:1, desc: 'Institute-Programme-Department mapping'},
    {id: 2, institute_id: 0, program_id:1, department_id:0, active:1, desc: 'Institute-Programme-Department mapping'},
    {id: 3, institute_id: 0, program_id:1, department_id:1, active:1, desc: 'Institute-Programme-Department mapping'},
  ],
  subjects: [
    {id: 0, code: 'Digital Communication', name: 'DGCM', status:1, description: 'List of Subjects'},
    {id: 1, code: 'Mathematics', name: 'MTMH', status:1, description: 'List of Subjects'},
  ],
  syllabuses: [
    {id: 0, offer_id: 0, subject_id: 1, sub_offered_sem:'First', sub_type:'R', description: 'R=Regular E=Elective'},
    {id: 1, offer_id: 0, subject_id: 0, sub_offered_sem:'Second', sub_type:'R', description: 'R=Regular E=Elective'},
  ],
  facaultys: [
    {id:0, name: 'name1', surname: 'surname1', designation: 'professior', mail: 'papubhat@gmail.com', mobile: '9963111687', address: 'address1', status:1, joining_date:  new Date('2014-07-21T16:15:34.479Z'), offers_id: 0},
    {id:1, name: 'name2', surname: 'surname2', designation: 'Lecturer', mail: 'papubhat@gmail.com', mobile: '9963111687', address: 'address1', status:1, joining_date:  new Date('2014-07-21T16:15:34.479Z'), offers_id: 0},
  ],
  teaches: [
    {id: 0, facaulty_id: 0, syllabus_id: 0, type: 0, description: '0=Lecture 1=Lab'},
    {id: 1, facaulty_id: 0, syllabus_id: 1, type: 0, description: '0=Lecture 1=Lab'},
    {id: 2, facaulty_id: 1, syllabus_id: 1, type: 0, description: '0=Lecture 1=Lab'},
  ],
  attendancemasters: [
    {id: 0, teach_id: 0, academiccalender_id: 0, batchno: 1, division:'A', description: 'Batch number for labs'},
    {id: 1, teach_id: 2, academiccalender_id: 0, batchno: 2, division:'A', description: 'Batch number for labs'},
    {id:2, teach_id: 1, academiccalender_id: 0, batchno: 2, division:'B', description: 'Batch number for labs'},
  ],
  students: [
    {id: 0, enrollmentno: '1821', rollno: '504321' , name: 'papu', father_name: 'promode', surname: 'bhat', email: 'pap@gmail.com', contact: '9963111687', parent_contact: '5636366492', address: 'hshhd', city: 'hyderabad', sem : 'First', status: 'C', offer_id: 0, div: 'A', batchno: 1, desc: "D=Detain L=Left C=Continue A=Alumni"},
    {id: 1, enrollmentno: '1822', rollno: '504322', name: 'papu1', father_name: 'promode', surname: 'bhat1', email: 'pap@gmail.com', contact: '9963111687', parent_contact: '5636366492', address: 'hshhd11', city: 'hyderabad', sem : 'First', status: 'C', offer_id: 0, div: 'A', batchno: 1, desc: "D=Detain L=Left C=Continue A=Alumni"}
  ],
  lectures: [
    {id: 0, date:  new Date(), attendancemaster_id: 0, active: 1, last_updated_on: new Date(), last_updated_by: 'papu'},
    {id: 1, date:  new Date(), attendancemaster_id: 1, active: 1, last_updated_on: new Date(), last_updated_by: 'papu'},
  ],
  attendence: [
    {id: 0, students_id:0, lectures_id: 0, presence: 1},
    {id: 1, students_id:0, lectures_id: 1, presence: 0},
  ],
};
