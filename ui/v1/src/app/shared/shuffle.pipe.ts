import { Pipe, PipeTransform } from '@angular/core';
import { isArray } from 'util';

@Pipe({
  name: 'shuffle'
})
export class ShufflePipe implements PipeTransform {

  transform(value: any[], args?: any): any[] {
    return this.shuffleArray(value);
  }

  private shuffleArray(array: any[]) {
    if (!isArray(array)) {
      return array;
    }

    const copy = [...array];

    for (let i = copy.length; i; --i) {
      const j = Math.floor(Math.random() * i);
      const x = copy[i - 1];
      copy[i - 1] = copy[j];
      copy[j] = x;
    }

    return copy;
  }

}
