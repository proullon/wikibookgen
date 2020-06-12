import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { HomeComponent } from './home/home.component';
import { ListWikibookComponent } from './list-wikibook/list-wikibook.component';
import { ShowWikibookComponent } from './show-wikibook/show-wikibook.component';
import { OrderWikibookComponent } from './order-wikibook/order-wikibook.component';
import { OrderStatusComponent } from './order-status/order-status.component';

const routes: Routes = [
  { path: '', redirectTo: '/home', pathMatch: 'full'},
  { path: 'home', component: HomeComponent},
  { path: 'wikibook/list', component: ListWikibookComponent},
  { path: 'wikibook/order', component: OrderWikibookComponent},
  { path: 'wikibook/order/:id', component: OrderStatusComponent},
  { path: 'wikibook/:id', component: ShowWikibookComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
