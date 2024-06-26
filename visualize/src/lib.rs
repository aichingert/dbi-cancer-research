use leptos::{html::Input, *};
use leptos_dom::log;
use serde::{Deserialize, Serialize};

const BEINGS: [&str; 3] = ["HUMAN", "MOUSE", "YEAST"];
const ENTER_KEY: u32 = 13;

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
    gene_id: i64,
    being_id: usize,
    name: String,
    essential_score: Option<f32>,
}

async fn fetch_lethal_genes<T>(gene_ident: String) -> Option<T>
where
    T: Serializable + std::fmt::Debug,
{
    let json = reqwasm::http::Request::get(&format!("http://127.0.0.1:3000/{}", gene_ident))
        .send()
        .await
        .map_err(|e| log!("{e}"))
        .ok()?
        .text()
        .await
        .ok()?;

    

    T::de(&json).ok()
}

#[component]
pub fn App() -> impl IntoView {
    let (gene, set_gene) = create_signal::<String>(String::new());
    let input_ref = create_node_ref::<Input>();
    let (_, set_pending) = create_signal(false);

    let get_gene = move || gene.get().clone();

    let lethal_genes = create_resource(
        get_gene,
        move |gene| async move { fetch_lethal_genes::<LethalGenes>(gene).await },
    );

    let find_gene = move |ev: web_sys::KeyboardEvent| {
        let input = input_ref.get().unwrap();
        ev.stop_propagation();

        if ev.key_code() == ENTER_KEY {
            let gene_ident = input.value().trim().to_string();

            if !gene_ident.is_empty() {
                set_gene.update(|gene| *gene = gene_ident);
            }
        }
    };

    view! {
        <h1>Synthetic Lethality</h1>

        <div class="input">
            <input
                placeholder="gene"
                on:keydown=find_gene
                node_ref=input_ref
            />
        </div>

        <Transition
            fallback=move || view! { <p>"Loading..."</p> }
            set_pending
        >
            {move || match lethal_genes.get() {
                    None => None,
                    Some(None) => Some(view! { <p> "No gene found" </p> }.into_any()),
                    Some(Some(genes)) => {
                        let mut score = genes.request_gene.essential_score.unwrap_or_default();

                        let pred = |g: &Lethal| -> f32 {
                            g.lethality_score * g.gene.essential_score.unwrap_or_default()
                        };

                        let lethality_score =
                            genes.yeast_genes.iter().map(pred).sum::<f32>()
                            + genes.human_genes.iter().map(pred).sum::<f32>()
                            + genes.mouse_genes.iter().map(pred).sum::<f32>();

                        score += (1. - score) * lethality_score;

                        Some(view! {
                            <div>
                                <h2>Checking Gene: </h2>
                                <h3>Lethality: {score}</h3>
                                <Gene gene=genes.request_gene.clone() />
                                <div class="row">
                                    <LethalGeneList lethal_genes=genes.human_genes.clone() />
                                    <LethalGeneList lethal_genes=genes.mouse_genes.clone() />
                                    <LethalGeneList lethal_genes=genes.yeast_genes.clone() />
                                </div>
                            </div>
                        }.into_any())
                    }
                }
            }
        </Transition>
    }
}

#[component]
pub fn Gene(gene: Gene) -> impl IntoView {
    view! {
        <div>
            <h3>Gene: </h3>
            <p>id: {gene.gene_id} </p>
            <p>name: {gene.name} </p>
            <p>being: {BEINGS[gene.being_id - 1]} </p>
            <p>essentiality:
            {move || if let Some(essential_score) = gene.essential_score {
                    essential_score.to_string()
                } else {
                    "Not essential".to_string()
                }
            }
            </p>
        </div>
    }
    .into_view()
}

#[component]
pub fn LethalGeneList(lethal_genes: Vec<Lethal>) -> impl IntoView {
    view! {
        <div class="column">
            <ul>
                <For
                    each=move || lethal_genes.clone()
                    key=|gene| gene.gene.gene_id
                    let:gene
                >
                    <Gene gene=gene.gene />
                    <p>Lethal score: {gene.lethality_score} </p>
                </For>
            </ul>
        </div>
    }
}
