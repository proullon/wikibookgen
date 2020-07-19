import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { NestedTreeControl } from '@angular/cdk/tree';
import { MatTreeModule, MatTreeNestedDataSource } from '@angular/material/tree';

import { WikibookgenService } from '../wikibookgen.service'; 

@Component({
  selector: 'app-order-status',
  templateUrl: './order-status.component.html',
  styleUrls: ['./order-status.component.sass']
})
export class OrderStatusComponent implements OnInit {

  orderStatus: string;

  constructor(
    private route: ActivatedRoute,
    private wikibookgenService: WikibookgenService,
    private location: Location
  ) { 
  }

  ngOnInit(): void {
    const uuid = this.route.snapshot.paramMap.get('id');
    this.getStatus(uuid);
  }

  getStatus(uuid: string): void {
    this.wikibookgenService.getOrderStatus(uuid) 
      .subscribe((orderStatus:string) => {
        console.log(orderStatus);
        this.orderStatus = orderStatus;
      });
  }

}
