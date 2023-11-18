use leptos::{error::Result, *};
use leptos_dom::log;
use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct LethalGenes {
    request_gene: Gene,
    human_genes: Vec<Lethal>,
    mouse_genes: Vec<Lethal>,
    yeast_genes: Vec<Lethal>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Lethal {
    gene: Gene,
    lethality_score: f32,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Gene {
    id: i64,
    being: i64,
    name: String,
    essentiality_score: Option<f32>,
}

#[derive(Error, Clone, Debug)]
pub enum GeneErro {
    #[error("No such human gene in our database")]
    InvalidGene
}

type Ident = String;

async fn fetch_lethal_genes<T>(gene_ident: String) -> Option<T> 
where
    T: Serializable + std::fmt::Debug,
{
    log!("{:?}", gene_ident);
    let json = reqwasm::http::Request::get(&format!("http://127.0.0.1:3000/{}", gene_ident))
        .send()
        .await
        .map_err(|e| log!("{e}"))
        .ok()?
        .text()
        .await
        .ok()?;

    let res = T::de(&json).ok();

    log!("{:?}", res);

    res
}

#[component]
pub fn App() -> impl IntoView {
    let (gene, set_gene) = create_signal::<Ident>(String::from("pten"));
    let (pending, set_pending) = create_signal(false);

    let get_gene = move || {
        gene.get().clone()
    };

    let lethal_genes = create_resource(
        move || (get_gene()), 
        move |gene| async move {
            let path = format!("{}", gene);
            fetch_lethal_genes::<LethalGenes>(path).await
        }
    );

    view! {
        <Transition
            fallback=move || view! { <p>"Loading..."</p> }
            set_pending
        >
            {move || match lethal_genes.get() {
                    None => None,
                    Some(None) => Some(view! { <p> "Error loading gene" </p> }.into_any()),
                    Some(Some(genes)) => {
                        Some(view! {
                            <div>
                                <h3>Gene: </h3>
                                <p> {genes.request_gene.id} </p>
                                <p> {genes.request_gene.name} </p>
                                <p> {genes.request_gene.being} </p>
                            </div>
                        }.into_any())
                    }
                }
            }
        </Transition>
    }
}
