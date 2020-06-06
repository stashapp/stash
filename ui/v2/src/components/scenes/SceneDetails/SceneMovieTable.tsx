import * as React from "react";
import { HTMLSelect, Divider} from "@blueprintjs/core";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";

type ValidTypes = GQL.SlimMovieDataFragment;

 export interface IProps {
   initialIds: string[] | undefined;
   initialIdx: string[] | undefined;
   onUpdate: (itemsNumber: string[]) => void; 
 }
let items: ValidTypes[];
let itemsFilter: ValidTypes[];
let storeIdx: string[];

export const SceneMovieTable: React.FunctionComponent<IProps> = (props: IProps) => {
const [itemsNumber, setItemsNumber] = React.useState<string[]>([]); 
const [initialIdsprev, setinitialIdsprev] = React.useState(props.initialIds);
const { data } = StashService.useAllMoviesForFilter(); 

items = !!data && !!data.allMovies ? data.allMovies : []; 
itemsFilter=[];
storeIdx=[];


if (!!props.initialIds && !!items && !!props.initialIdx) 
{
   for(var i=0; i< props.initialIds!.length; i++)
   {
      itemsFilter=itemsFilter.concat(items.filter((x) => x.id ===props.initialIds![i])); 
    }
 
}
   /* eslint-disable react-hooks/rules-of-hooks */
   React.useEffect(() => {
     if (!!props.initialIdx) 
     {
        setItemsNumber(props.initialIdx);
       }
     }, [props.initialIdx]);
   /* eslint-enable */
   

    React.useEffect(() => {
      if (!!props.initialIds) {
        setinitialIdsprev(props.initialIds);
        UpdateIndex();
      }
  }, [props.initialIds]);

   const updateFieldChanged = (index : any)  => (e : any) => {
      let newArr = [...itemsNumber]; 
      newArr[index] = e.target.value; 
      setItemsNumber(newArr); 
      props.onUpdate(newArr);
     }
  
    const updateIdsChanged = (index : any, value: string) => {
      storeIdx.push(value); 
      setItemsNumber(storeIdx); 
      props.onUpdate(storeIdx);
     }


function UpdateIndex(){

 if (!!props.initialIds && !!initialIdsprev ){
 loop1:
  for(var i=0; i< props.initialIds!.length; i++) { 
    for(var j=0; j< initialIdsprev!.length; j++) {
      
        if (props.initialIds[i]===initialIdsprev[j])
        {
            updateIdsChanged(i, props.initialIdx![j]);
            continue loop1;
        } 
    }
      updateIdsChanged(i, "0");
     }
 
 }
}

   function renderTableData() {

     return(
        <tbody>
            { itemsFilter!.map((item, index : any)  => ( 
               <tr>
                  <td>{item.name} </td>
                  <td><Divider /> </td>
                  <td key={item.toString()}> Scene number: <HTMLSelect
                        options={["","1", "2", "3", "4", "5","6","7","8","9","10"]}
                        onChange={updateFieldChanged(index)}
                        value={itemsNumber[index]}
                        />
                  </td>
               </tr>
         ))}    
        </tbody>
      )
     }

     return (

        <div>
          <table id='movies'>
              {renderTableData()}
          </table>
        </div>
        );
};




