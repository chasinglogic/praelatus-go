<template>
    <b-navbar fixed=top toggle-breakpoint="md" toggleable
      :class="bg-praelatus">
        <b-nav-toggle target="nav_collapse"></b-nav-toggle>

        <b-collapse is-nav id="nav_collapse">
            <b-navbar-brand href="/">
                <img src="/assets/static/img/logo_arrow.svg"
                     height="30"
                     width="30">
            </b-navbar-brand>
            <b-nav is-nav-bar>
                <b-nav-item href="/projects/search">Projects</b-nav-item>
                <b-nav-item-dropdown text="Tickets">
                  <b-dropdown-item to="/queries/mine">Manage Queries</b-dropdown-item>
                  <b-dropdown-item to="/queries">Search</b-dropdown-item>
                </b-nav-item-dropdown>
            </b-nav>
            <b-nav is-nav-bar class="ml-auto">
              <b-nav-item class="btn btn-success" 
                href="/tickets/create">Create</b-nav-item>
                <li class="nav-item dropdown btn-group">
                    <a class="dropdown-toggle nav-link"
                       id="userMenu1"
                       data-toggle="dropdown"
                       aria-haspopup="true" aria-expanded="false">
                        <img width="30" height="30" src="">
                    </a>
                    <div class="dropdown-menu dropdown-menu-right"
                         aria-labelledby="userMenu1">
                        <a class="dropdown-item" href="/admin/">
                            System Administration
                        </a>
                        <a class="dropdown-item" href="/logout">
                            Log Out
                        </a>
                    </div>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="">
                        Sign In
                    </a>
                </li>
                <!-- Space out the buttons -->
                <li class="nav-item disabled hidden-sm-down">
                    <a class="nav-link" href="#">/</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link"
                       href="">
                        Sign Up
                    </a>
                </li>
            </b-nav>
        </b-collapse>
    </b-navbar>
</template>
