import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { StashService } from '../../core/stash.service';

import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

import { FormControl } from '@angular/forms';

@Component({
  selector: 'app-performer-form',
  templateUrl: './performer-form.component.html',
  styleUrls: ['./performer-form.component.css']
})
export class PerformerFormComponent implements OnInit, OnDestroy {
  name: string;
  favorite: boolean;
  aliases: string;
  country: string;
  birthdate: string;
  ethnicity: string;
  eye_color: string;
  height: string;
  measurements: string;
  fake_tits: string;
  career_length: string;
  tattoos: string;
  piercings: string;
  url: string;
  twitter: string;
  instagram: string;
  image: string;

  loading = true;
  imagePreview: string;
  image_path: string;
  ethnicityOptions: string[] = ['white', 'black', 'asian', 'hispanic'];
  performerNameOptions: string[] = [];
  selectedName: string;

  searchFormControl = new FormControl();

  constructor(
    private route: ActivatedRoute,
    private stashService: StashService,
    private router: Router
  ) {}

  ngOnInit() {
    this.getPerformer();

    this.searchFormControl.valueChanges.pipe(
                            debounceTime(400),
                            distinctUntilChanged()
                          ).subscribe(term => {
                            this.name = term;
                            this.getPerformerNames(term);
                          });
  }

  ngOnDestroy() {}

  async getPerformer() {
    const id = parseInt(this.route.snapshot.params['id'], 10);
    if (!!id === false) {
      console.log('new performer');
      this.loading = false;
      return;
    }

    const result = await this.stashService.findPerformer(id).result();
    this.loading = result.loading;

    this.name = result.data.findPerformer.name;
    this.selectedName = this.name;
    this.searchFormControl.setValue(this.name);
    this.favorite = result.data.findPerformer.favorite;
    this.aliases = result.data.findPerformer.aliases;
    this.country = result.data.findPerformer.country;
    this.birthdate = result.data.findPerformer.birthdate;
    this.ethnicity = result.data.findPerformer.ethnicity;
    this.eye_color = result.data.findPerformer.eye_color;
    this.height = result.data.findPerformer.height;
    this.measurements = result.data.findPerformer.measurements;
    this.fake_tits = result.data.findPerformer.fake_tits;
    this.career_length = result.data.findPerformer.career_length;
    this.tattoos = result.data.findPerformer.tattoos;
    this.piercings = result.data.findPerformer.piercings;
    this.url = result.data.findPerformer.url;
    this.twitter = result.data.findPerformer.twitter;
    this.instagram = result.data.findPerformer.instagram;

    this.image_path = result.data.findPerformer.image_path;
    this.imagePreview = this.image_path;
  }

  async getPerformerNames(query: string) {
    if (query === undefined) { return; }
    if (this.selectedName !== this.name) { this.selectedName = null; }
    const result = await this.stashService.scrapeFreeonesPerformers(query).result();
    this.performerNameOptions = result.data.scrapeFreeonesPerformerList;
  }

  onClickedPerformerName(name) {
    this.name = name;
    this.selectedName = name;
    this.searchFormControl.setValue(this.name);
  }

  onImageChange(event) {
    const file: File = event.target.files[0];
    const reader: FileReader = new FileReader();

    reader.onloadend = (e) => {
      this.image = reader.result as string;
      this.imagePreview = this.image;
    };
    reader.readAsDataURL(file);
  }

  onResetImage(imageInput) {
    imageInput.value = '';
    this.imagePreview = this.image_path;
    this.image = null;
  }

  onFavoriteChange() {
    this.favorite = !this.favorite;
  }

  onSubmit() {
    const id = this.route.snapshot.params['id'];

    if (!!id) {
      this.stashService.performerUpdate({
        id: id,
        name: this.name,
        url: this.url,
        birthdate: this.birthdate,
        ethnicity: this.ethnicity,
        country: this.country,
        eye_color: this.eye_color,
        height: this.height,
        measurements: this.measurements,
        fake_tits: this.fake_tits,
        career_length: this.career_length,
        tattoos: this.tattoos,
        piercings: this.piercings,
        aliases: this.aliases,
        twitter: this.twitter,
        instagram: this.instagram,
        favorite: this.favorite,
        image: this.image
      }).subscribe(result => {
        this.router.navigate(['/performers', id]);
      });
    } else {
      this.stashService.performerCreate({
        name: this.name,
        url: this.url,
        birthdate: this.birthdate,
        ethnicity: this.ethnicity,
        country: this.country,
        eye_color: this.eye_color,
        height: this.height,
        measurements: this.measurements,
        fake_tits: this.fake_tits,
        career_length: this.career_length,
        tattoos: this.tattoos,
        piercings: this.piercings,
        aliases: this.aliases,
        twitter: this.twitter,
        instagram: this.instagram,
        favorite: this.favorite,
        image: this.image
      }).subscribe(result => {
        this.router.navigate(['/performers', result.data.performerCreate.id]);
      });
    }
  }

  async onScrape() {
    this.loading = true;
    const result = await this.stashService.scrapeFreeones(this.name).result();
    this.loading = false;

    this.url           = result.data.scrapeFreeones.url;
    this.name          = result.data.scrapeFreeones.name;
    this.searchFormControl.setValue(this.name);
    this.aliases       = result.data.scrapeFreeones.aliases;
    this.country       = result.data.scrapeFreeones.country;
    this.birthdate     = result.data.scrapeFreeones.birthdate ? result.data.scrapeFreeones.birthdate : this.birthdate;
    this.ethnicity     = result.data.scrapeFreeones.ethnicity;
    this.eye_color     = result.data.scrapeFreeones.eye_color;
    this.height        = result.data.scrapeFreeones.height;
    this.measurements  = result.data.scrapeFreeones.measurements;
    this.fake_tits     = result.data.scrapeFreeones.fake_tits;
    this.career_length = result.data.scrapeFreeones.career_length;
    this.tattoos       = result.data.scrapeFreeones.tattoos;
    this.piercings     = result.data.scrapeFreeones.piercings;
    this.twitter       = result.data.scrapeFreeones.twitter;
    this.instagram     = result.data.scrapeFreeones.instagram;
  }
}
