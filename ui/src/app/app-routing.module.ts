import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { ListWikibookComponent } from './list-wikibook/list-wikibook.component';

const routes: Routes = [
//  { path: '', redirectTo: '/home', pathMatch: 'full'},
  { path: 'wikibook/list', component: ListWikibookComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
