import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'age'
})
export class AgePipe implements PipeTransform {

  transform(value: string, ageFromDate?: string): number {
    if (!!value === false) { return 0; }

    const birthdate = new Date(value);
    const fromDate = !!ageFromDate ? new Date(ageFromDate) : new Date();

    let age = fromDate.getFullYear() - birthdate.getFullYear();
    if (birthdate.getMonth() > fromDate.getMonth() ||
        (birthdate.getMonth() >= fromDate.getMonth() && birthdate.getDay() > fromDate.getDay())) {
      age -= 1;
    }

    return age;
  }

}
