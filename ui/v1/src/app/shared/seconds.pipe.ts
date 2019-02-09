import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'seconds'
})
export class SecondsPipe implements PipeTransform {

  transform(value: any, args?: any): any {
    return new Date(value * 1000).toISOString().substr(11, 8);
  }

}
