import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'filename'
})
export class FileNamePipe implements PipeTransform {

  transform(value: string, args?: any): string {
    if (!!value === false) { return 'No File Name'; }
    return value.replace(/^.*[\\\/]/, '');
  }

}
