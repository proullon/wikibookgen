import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { HomeComponent } from './home/home.component';
import { ListWikibookComponent } from './list-wikibook/list-wikibook.component';
import { OrderWikibookComponent } from './order-wikibook/order-wikibook.component';

const routes: Routes = [
  { path: '', redirectTo: '/home', pathMatch: 'full'},
  { path: 'home', component: HomeComponent},
  { path: 'wikibook/list', component: ListWikibookComponent},
  { path: 'wikibook/order', component: OrderWikibookComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
