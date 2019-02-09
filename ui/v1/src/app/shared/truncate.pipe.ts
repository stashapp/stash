import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'truncate'
})
export class TruncatePipe implements PipeTransform {

  transform(value: string, limit: number = 100, tail: string = '...'): string {
    return value.length > limit ? value.substring(0, limit) + tail : value;
  }

}
