import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { Observable, of } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';

import { Wikibook, ListWikibookRequest, ListWikibookResponse, CompleteRequest, OrderRequest } from './wikibook';
import { environment } from './../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class WikibookgenService {

  private api = environment.api;

  constructor(
    private http: HttpClient,
  ) { }

  public listWikibook(page: number, size: number, language: string): Observable<Wikibook[]> {
    console.log('listWikibook ' + page + ' ' + size);
    return this.http.get<Wikibook[]>(`${this.api}/wikibook?page=${page}&size=${size}&language=${language}`)
      .pipe(
        map((result:any)=>{
          return result.wikibooks
        }),
        catchError(this.handleError('listWikibook', null))
      );
  }
  
  public getWikibook(uuid: string): Observable<Wikibook> {
    console.log('getWikibook ' + uuid);
    return this.http.get<Wikibook>(`${this.api}/wikibook/${uuid}`)
      .pipe(
        map((result:any)=>{
          return result.wikibook
        }),
        catchError(this.handleError('getWikibook', null))
      );
  }

  public autocomplete(value: string, language: string): Observable<string[]> {
    console.log('complete ' + value);
    var req: CompleteRequest = {
      value: value,
      language: language
    }
    return this.http.post<string[]>(`${this.api}/complete`, req)
      .pipe(
        map((result:any)=>{
          return result.titles
        }),
        catchError(this.handleError('complete', null))
      );
  }
  
  public order(subject: string, language: string, model: string): Observable<string> {
    console.log('ordering '+ model + ' ' + subject + ' in ' + language);
    var req: OrderRequest = {
      subject: subject,
      language: language,
      model: model,
    }
    return this.http.post<string[]>(`${this.api}/order`, req)
      .pipe(
        map((result:any)=>{
          return result.uuid
        }),
        catchError(this.handleError('complete', null))
      );
  }

  /** handle failed HTTP operation
   *  let the app continue
   *  @param operation - name of the failed operation
   *  @param result - optional value to return as the observable result
   **/
  private handleError<T> (operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.log(`${operation} failed: ${error}`);
      return of(result as T);
    };
  }
}
