# `flack`

`flack` is a thin wrapper around your usual nix commands `nix profile add/upgrade` and `nix flake update`, but centered around one specific workflow: grouping the nix-packaged tools you often use inside "toolkits", installing these toolkits instead of individual tools in your Nix profile, and managing your whole profile as a _stack_.

This results in a coarser-grained and more manageable profile, in which you can directly identify your most fundamental tools (at the bottom of the stack, meant to stay) and the more "transient" ones (at the top) possibly even specific to some project.

`nix profile` has a little-known feature which is that every installed package is given a **priority** value, settable at install time.
If you install two packages which provide the same executable, the package with the lowest priority value will take precedence, and the executable from the other package will simply be hidden.

_NOTE: All em dashes in the following man-made slop are man-made too, painstakingly input one by one._

_NOTE2: ...the code, on the other hand, *is* vibe-coded with GLM-5.1 like there is no tomorrow._

## Usage

Install `flack` using `flack`:

```sh
nix run github:YPares/flack.nix -- push github:YPares/flack.nix#flack
```

then `flack` shows the current contents of your stac... I mean your profile.

Then use:

- `flack push <flake>#<package>` just as you would use `nix profile add`
- `flack pop` to remove the last installed package
- `flack inputs [--flake <path>]` to list info about a flake's inputs (full URL, revision, last update date)
- `flack update [--flake <path>]` to interactively select which inputs of a flake to update in its lockfile
- `flack` upgrade` to interactively select which elements of your profile to upgrade

## But doesn't it completely kill declarativity?

Short answer: not _completely_.

Long answer:

To me, keeping the advantages of a declarative setup (i.e. one file or repository with your entire config — or possibly several different configs, one _"apply this configuration now whatever the current state"_ unique command) is not just a matter of tooling, it's a matter of _granularity_.

If you were to manually install every single package this way, it would very much kill the benefits of declarativity indeed. However, you can use this approach to power a semi-declarative, yet semi-valuable workflow if paired with custom "environments" built via the nix function `buildEnv`, which does the aforementioned grouping.

Declarative management of your setup, through complete `NixOS` or `home-manager` configurations, is awesome for servers or if you really need one specific environment which you want to replicate on various machines, but I believe there is a place for a compromise between this and purely imperative, `apt-get`-like package management.

## Why not just `nix shell/develop` or `nix run`?

All the other way across the "declarativity spectrum", you have the ability to define a local shell with just some specific packages — or even download and run a Nix package — in both cases possibly locked to a very specific version if you run from a specific flake which you keep pinned. Excellent for just running one-off commands, but that comes with some caveats:

- nix devshells are difficult to compose. They are just not meant for that. You can start shells within shells but that becomes awkward. `buildEnv` naturally composes both declaratively (you can reuse some env as the input of another `buildEnv` call) and imperatively (you can install in your profile as many envs as you want, coming from different flakes, each possibly using different versions of nixpkgs)
- `nix run nixpkgs#something` — besides being more cumbersome to type — means you forfeit locking entirely (_unless_ you pin flakes at the Nix registry level, which is another can of worms and comes with its own problems)

## This is just plain Nix, nothing fancy

Good thing is that the Nix ecosystem already has everything needed to power such a workflow:

- the `buildEnv` function to define such toolkits
- nix flakes to lock and upgrade at will the packages contained in the toolkits
- nix profiles to install and remove them at will, with the notion of priority to deal with conflicts
- `blueprint` or `flake-utils` to easily define your various envs (toolkits) as output packages of a flake

So don't think of `flack` as a new competitor to `niv`, `home-manager`, `devenv` or similar solutions, just as a light enabler for a specific workflow which I've been using for several years and which has provided me great value due to its flexibility, while keeping all of the advantages of Nix intact.

That's the beauty of Nix to me: from a common language and package sources, you can have very different workflows built on top, depending on what suits your needs best. If the single "TRVE DEKLÄRATIV" approach was for everybody and every situation, we'd have pinned the whole internet by now.
