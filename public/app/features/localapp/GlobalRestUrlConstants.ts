export class GlobalRestUrlConstants {
  PROTOCOL: string;
  IP: string;
  PORT: string;
  ROOT_API: string;
  BASE_URL: string;
  constructor() {
    this.PROTOCOL = 'http';
    this.IP = 'localhost'; //18.209.4.2
    this.PORT = '8080';
    this.ROOT_API = 'api';
    this.BASE_URL = this.PROTOCOL + '://' + this.IP + ':' + this.PORT + '/' + this.ROOT_API;
  }
  getBaseRestUrl() {
    return this.BASE_URL;
  }
  // getCollegeRestUrl() {
  //   return this.getBaseRestUrl() + '/cmscollege/';
  // }
  getStateRestUrl() {
    return this.getBaseRestUrl() + '/states/';
  }
  getCityRestUrl() {
    return this.getBaseRestUrl() + '/cities/';
  }
  getBranchRestUrl() {
    return this.getBaseRestUrl() + '/cmsbranches/';
  }
  getBranchesByCollegeIdRestUrl() {
    return this.getBaseRestUrl() + '/cmsbranches-collegeid/';
  }
  getLegalEntitiesRestUrl() {
    return this.getBaseRestUrl() + '/cmslegal-entities/';
  }
  getAuthorizedSignatoryRestUrl() {
    return this.getBaseRestUrl() + '/cmsauthorized-signatories/';
  }
  getAuthorizedSignatoriesByCollegeIdRestUrl() {
    return this.getBaseRestUrl() + '/cmsauthorized-signatories-collegeid/';
  }
  getBankAccountRestUrl() {
    return this.getBaseRestUrl() + '/cmsbank-accounts/';
  }
  getBankAccountsByCollegeIdRestUrl() {
    return this.getBaseRestUrl() + '/cmsbank-accounts-collegeid/';
  }
  getAcademicYearRestUrl() {
    return this.getBaseRestUrl() + '/cmsacademic-years/';
  }
  getDepartmentRestUrl() {
    return this.getBaseRestUrl() + '/cmsdepartments/';
  }
  getDepartmentByBranchIdRestUrl() {
    return this.getBaseRestUrl() + '/cmsdepartments-branchid/';
  }
  getBatchRestUrl() {
    return this.getBaseRestUrl() + '/cmsbatches/';
  }
  getBatchByDepartmentIdRestUrl() {
    return this.getBaseRestUrl() + '/cmsbatches-departmentid/';
  }
  getTeacherRestUrl() {
    return this.getBaseRestUrl() + '/cmsteachers/';
  }
  getTeacherByQueryParamsRestUrl() {
    return this.getBaseRestUrl() + '/cmsteachers-qryprms?';
  }
  getHolidayRestUrl() {
    return this.getBaseRestUrl() + '/cmsholidays/';
  }
  getHolidayByAcademicYearIdRestUrl() {
    return this.getBaseRestUrl() + '/cmsholidays-by_academicyearid';
  }
  getTermRestUrl() {
    return this.getBaseRestUrl() + '/cmsterms/';
  }
  getTermByAcademicYearIdRestUrl() {
    return this.getBaseRestUrl() + '/cmsterms-by_academicyearid';
  }
  getSubjectRestUrl() {
    return this.getBaseRestUrl() + '/cmssubjects/';
  }

  getSubjectByDeptBatchIdRestUrl() {
    return this.getBaseRestUrl() + '/cmssubjects?';
  }

  getSemesterRestUrl() {
    return this.getBaseRestUrl() + '/cmssemesters/';
  }
  getSectionRestUrl() {
    return this.getBaseRestUrl() + '/cmssections/';
  }
  getSectionByBatchRestUrl() {
    return this.getBaseRestUrl() + '/cmssections-batchid/';
  }

  getCmsLecturesUrl() {
    return this.getBaseRestUrl() + '/cmslectures?';
  }

  getAttendanceMasterUrl() {
    return this.getBaseRestUrl() + '/cmsattendance-masters';
  }

  getAttendanceMasterByBatchAndSectioinUrl() {
    return this.getBaseRestUrl() + '/cmsattendance-masters-bybatchsection';
  }

  getSelectedLectures() {
    return this.getBaseRestUrl() + '/cmsmeta-lecture-selected';
  }

  getAllLecturesReport() {
    return this.getBaseRestUrl() + '/cmsmeta-lecture-viewalllectures';
  }

  getUploadCmsDataUrl() {
    return this.getBaseRestUrl() + '/cmsdataimport';
  }
}
